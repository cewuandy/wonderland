package usecase

import (
	"context"
	"fmt"
	"github.com/samber/do"
	orderedmap "github.com/wk8/go-ordered-map/v2"
	"math"
	"strconv"
	"strings"

	"github.com/cewuandy/wonderland/internal/domain"
)

type materialUseCase struct {
	fileRepo domain.FileRepo

	materialMap map[string]*domain.Material
}

func (m *materialUseCase) GetProductionProcess(ctx context.Context,
	req *domain.ProductionProcessReq) (string, error) {
	var (
		result          string
		tree            string
		productionLevel string
		originMaterial  string
		totalTimeStr    string
		totalTime       float32
		countResult     map[string]float32
		orderedResults  []*orderedmap.OrderedMap[string, string]
		err             error
	)
	tree, countResult, err = m.getTree(ctx, req)
	if err != nil {
		return "", err
	}

	productionLevel, orderedResults, totalTime, err = m.getProductionLevel(ctx, req, countResult)
	if err != nil {
		return "", err
	}

	originMaterial = m.getOriginMaterial(orderedResults)
	totalTimeStr = m.getTotalTime(totalTime)

	result = fmt.Sprintf("%s%s%s%s", tree, productionLevel, originMaterial, totalTimeStr)

	return result, nil
}

func (m *materialUseCase) getTree(ctx context.Context, req *domain.ProductionProcessReq) (string,
	map[string]float32, error) {
	var (
		bias        = m.getBias(req)
		result      string
		countResult map[string]float32
	)

	material, ok := m.materialMap[req.Name]
	if !ok {
		return "", nil, fmt.Errorf("item is not found")
	}

	reqQuantity := float32(req.Quantity)
	result += fmt.Sprintf("%s (%s)\n", material.Name, material.Tool)

	result, countResult = m.search(
		req.Name, &result, reqQuantity, bias, 1, map[string]float32{},
	)
	return result, countResult, nil
}

func (m *materialUseCase) getProductionLevel(ctx context.Context,
	req *domain.ProductionProcessReq, countResult map[string]float32) (string,
	[]*orderedmap.OrderedMap[string, string], float32, error) {
	var (
		orderedResults = make([]*orderedmap.OrderedMap[string, string], 10)
		bias           = m.getBias(req)
		totalTime      float32
	)

	{
		material, _ := m.materialMap[req.Name]
		totalTime += material.Time * float32(req.Quantity)
		countResult[req.Name] = float32(req.Quantity)
	}

	for k, v := range countResult {
		level := m.getLevel(k, 0, 0)
		if orderedResults[level] == nil {
			orderedResults[level] = orderedmap.New[string, string]()
		}
		orderedResults[level].Set(k, fmt.Sprintf("x%.0f ", math.Ceil(float64(v))))
	}

	var result string
	for i := len(orderedResults) - 1; i > 0; i-- {
		if orderedResults[i].Len() == 0 {
			continue
		}

		if orderedResults[i-1].Len() > 5 && i > 1 {
			err := m.iterateAndMoveToFront(orderedResults[i], orderedResults[i-1])
			if err != nil {
				return "", nil, 0.0, err
			}
		}

		var levelTotalTime float32
		for pair := orderedResults[i].Newest(); pair != nil; pair = pair.Prev() {
			material, _ := m.materialMap[pair.Key]
			materialCount, _ := strconv.ParseFloat(pair.Value[1:len(pair.Value)-1], 32)
			materialTime := material.Time * float32(materialCount) * bias
			levelTotalTime += materialTime
			result = fmt.Sprintf(
				"%s%s(%s) %.1fm\n", pair.Key, pair.Value, material.Tool, materialTime,
			) + result
		}
		result = fmt.Sprintf("\nLevel: %d (%.1fm)\n", i, levelTotalTime) + result
		totalTime += levelTotalTime
	}

	result += "\n"

	return result, orderedResults, totalTime, nil
}

func (m *materialUseCase) getOriginMaterial(orderedResults []*orderedmap.OrderedMap[string, string]) string {
	var result string
	result += "原材料列表：\n"
	for pair := orderedResults[0].Oldest(); pair != nil; pair = pair.Next() {
		result += fmt.Sprintf("%s%s\n", pair.Key, pair.Value)
	}
	return result
}

func (m *materialUseCase) getTotalTime(totalTime float32) string {
	var result string
	hours := int(totalTime) / 60
	minutes := int(totalTime) - hours*60
	seconds := (totalTime - float32(hours*60) - float32(minutes)) * 60
	result += "\n"
	result += fmt.Sprintf("總耗時：%d時 %d分 %.0f秒\n", hours, minutes, seconds)
	return result
}

func (m *materialUseCase) initMaterialMap() {
	var (
		unfinishedQueue []*domain.Material
		fileContentMap  = map[string][]*string{}
		materialMap     = map[string]*domain.Material{}
		err             error
	)

	unfinishedQueue = []*domain.Material{}

	fileContentMap, err = m.fileRepo.ListFilesContent()
	if err != nil {
		panic(err)
	}

	for k := range fileContentMap {
		lines := fileContentMap[k]
		for i := 0; i < len(lines); i++ {
			titleAndMaterial := strings.Split(*lines[i], "=")
			titleList := strings.Split(titleAndMaterial[0], "x")
			title := titleList[0]
			time, _ := strconv.ParseFloat(titleList[2], 32)
			count, _ := strconv.ParseFloat(titleList[1], 32)
			quantity, _ := strconv.ParseInt(titleList[1], 10, 8)
			titleQuantity := float32(quantity)
			titleTime := time / count

			materialMap[title] = &domain.Material{
				Name:     title,
				Quantity: 1,
				Tool:     strings.TrimRight(k, ".txt"),
				Time:     float32(titleTime),
				Dependencies: m.getDependencies(
					materialMap, titleQuantity, titleAndMaterial[1], &unfinishedQueue,
				),
			}
		}
	}

	for {
		if len(unfinishedQueue) == 0 {
			break
		}
		dep := unfinishedQueue[0]
		unfinishedQueue = unfinishedQueue[1:]
		material, ok := materialMap[dep.Name]
		if ok {
			dep.Tool = material.Tool
		}
	}

	m.materialMap = materialMap
}

func (m *materialUseCase) getDependencies(materialMap map[string]*domain.Material,
	titleQuantity float32, materialString string,
	unfinishedQueue *[]*domain.Material) []*domain.Material {
	var materials []*domain.Material
	materialList := strings.Split(materialString, " ")
	for _, mStr := range materialList {
		if len(mStr) == 0 {
			continue
		}
		// titleQuantityStr[0] equals item name, titleQuantityStr[1] equals item quantity
		titleQuantityStr := strings.Split(mStr, "x")
		itemName := titleQuantityStr[0]
		itemQuantity := titleQuantityStr[1]
		quantity, _ := strconv.ParseInt(itemQuantity, 10, 8)
		material, ok := materialMap[itemName]
		if ok {
			tool := material.Tool
			materials = append(
				materials, &domain.Material{
					Name:     itemName,
					Tool:     tool,
					Quantity: float32(quantity) / titleQuantity,
				},
			)
		} else {
			dep := &domain.Material{Name: itemName, Quantity: float32(quantity) / titleQuantity}
			materials = append(materials, dep)
			*unfinishedQueue = append(*unfinishedQueue, dep)
		}
	}

	return materials
}

func (m *materialUseCase) getBias(req *domain.ProductionProcessReq) float32 {
	var bias = float32(1)
	if req.ClockEnable {
		bias -= 0.02
	}
	if req.StandClockEnable {
		bias -= 0.03
	}
	if req.WindmillEnable {
		bias -= 0.03
	}
	if req.ACEnable {
		bias -= 0.14
	}

	return bias
}

func (m *materialUseCase) search(itemName string, result *string, count, bias float32, level int,
	countResult map[string]float32) (string, map[string]float32) {
	v, ok := m.materialMap[itemName]
	if !ok {
		return "", nil
	}

	for _, dep := range v.Dependencies {
		var timeConsumed = float32(0)
		material, ok := m.materialMap[dep.Name]
		if ok {
			timeConsumed = count * material.Time * dep.Quantity * bias
		}
		c := count * dep.Quantity
		countResult[dep.Name] += c

		for i := 0; i < level; i++ {
			*result += fmt.Sprint("  ")
		}
		*result += fmt.Sprintf("%sx%.0f ", dep.Name, c)
		if len(dep.Tool) == 0 {
			*result += "\n"
		} else {
			*result += fmt.Sprintf("(%s) %.1fm\n", dep.Tool, timeConsumed)
		}
		m.search(dep.Name, result, c, bias, level+1, countResult)
	}

	return *result, countResult
}

func (m *materialUseCase) getLevel(itemName string, level int, max int) int {
	v, ok := m.materialMap[itemName]
	if !ok {
		if max < level {
			max = level
		}
		return max
	}

	for _, dep := range v.Dependencies {
		max = m.getLevel(dep.Name, level+1, max)
	}

	return max
}

func (m *materialUseCase) iterateAndMoveToFront(currentLevel *orderedmap.OrderedMap[string, string],
	nextLevel *orderedmap.OrderedMap[string, string]) error {
	for pair := currentLevel.Oldest(); pair != nil; pair = pair.Next() {
		material, _ := m.materialMap[pair.Key]
		for _, d := range material.Dependencies {
			_, ok := nextLevel.Get(d.Name)
			if !ok {
				continue
			}
			err := nextLevel.MoveToFront(d.Name)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func NewMaterialUseCase(i *do.Injector) (domain.MaterialUseCase, error) {
	usecase := &materialUseCase{do.MustInvoke[domain.FileRepo](i), nil}
	usecase.initMaterialMap()
	return usecase, nil
}
