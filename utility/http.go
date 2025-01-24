package utility

import (
	"fmt"
	"strconv"
	"strings"
)

// parseFilterConditions 解析查询字符串为过滤条件。
//
// 参数:
//   - filter ([]string): 包含过滤条件的切片。
//
// 返回:
//   - ([]string, error): 解析后的过滤条件或解析失败时的错误。
func parseFilterConditions(filter []string) ([][]string, error) {
	switch filter[0] {
	case "equal", "eq":
		c, err := strconv.Atoi(filter[1])
		if err != nil {
			return nil, err
		}
		if c%2 != 0 {
			return nil, fmt.Errorf("参数数量错误")
		}

		var result [][]string
		for i := 0; i < c; i += 2 {
			result = append(result, []string{"equal", filter[2+i], filter[3+i]})
		}
		return result, nil
	case "not-equal", "ne":
		c, err := strconv.Atoi(filter[1])
		if err != nil {
			return nil, err
		}
		if c%2 != 0 {
			return nil, fmt.Errorf("参数数量错误")
		}

		var result [][]string
		for i := 0; i < c; i += 2 {
			result = append(result, []string{"not-equal", filter[2+i], filter[3+i]})
		}
		return result, nil
	case "in":
		c, err := strconv.Atoi(filter[1])
		if err != nil {
			return nil, err
		}
		v := filter[2 : 2+c]
		return [][]string{append([]string{"in"}, v...)}, nil
	case "like", "lk":
		c, err := strconv.Atoi(filter[1])
		if err != nil {
			return nil, err
		}
		if c%2 != 0 {
			return nil, fmt.Errorf("参数数量错误")
		}

		var result [][]string
		for i := 0; i < c; i += 2 {
			result = append(result, []string{"like", filter[2+i], filter[3+i]})
		}
		return result, nil
	case "greater-equal", "ge":
		c, err := strconv.Atoi(filter[1])
		if err != nil {
			return nil, err
		}
		if c%2 != 0 {
			return nil, fmt.Errorf("参数数量错误")
		}

		var result [][]string
		for i := 0; i < c; i += 2 {
			result = append(result, []string{"greater-equal", filter[2+i], filter[3+i]})
		}
		return result, nil
	case "less-equal", "le":
		c, err := strconv.Atoi(filter[1])
		if err != nil {
			return nil, err
		}
		if c%2 != 0 {
			return nil, fmt.Errorf("参数数量错误")
		}

		var result [][]string
		for i := 0; i < c; i += 2 {
			result = append(result, []string{"less-equal", filter[2+i], filter[3+i]})
		}
		return result, nil
	case "greater", "gt":
		c, err := strconv.Atoi(filter[1])
		if err != nil {
			return nil, err
		}
		if c%2 != 0 {
			return nil, fmt.Errorf("参数数量错误")
		}

		var result [][]string
		for i := 0; i < c; i += 2 {
			result = append(result, []string{"greater", filter[2+i], filter[3+i]})
		}
		return result, nil
	case "less", "lt":
		c, err := strconv.Atoi(filter[1])
		if err != nil {
			return nil, err
		}
		if c%2 != 0 {
			return nil, fmt.Errorf("参数数量错误")
		}

		var result [][]string
		for i := 0; i < c; i += 2 {
			result = append(result, []string{"less", filter[2+i], filter[3+i]})
		}
		return result, nil
	case "array-contain", "act":
		c, err := strconv.Atoi(filter[1])
		if err != nil {
			return nil, err
		}
		v := filter[2 : 2+c]
		return [][]string{append([]string{"json-array-contains"}, v...)}, nil
	case "object-contain", "oct":
		c, err := strconv.Atoi(filter[1])
		if err != nil {
			return nil, err
		}
		if c%3 != 0 {
			return nil, fmt.Errorf("参数数量错误")
		}

		var result [][]string
		for i := 0; i < c; i += 3 {
			result = append(result, []string{"json-object-contains", filter[2+i], filter[3+i], filter[4+i]})
		}
		return result, nil
	}
	return nil, nil
}

// ConvertQueryStringToDefaultFilter 将查询字符串解析为默认过滤器。
//
// 参数:
//   - qs (string): 原始查询字符串。
//
// 返回:
//   - ([][]string, error): 解析后的过滤条件切片或解析失败时的错误。
func ConvertQueryStringToDefaultFilter(qs string) ([][]string, error) {
	result := [][]string{}
	if qs == "" {
		return result, nil
	}
	filter := strings.Split(qs, ",")
	for len(filter) > 0 {
		qty, err := strconv.Atoi(filter[1])
		if err != nil {
			return nil, err
		}
		p := filter[0 : 2+qty]
		parameter, err := parseFilterConditions(p)
		if err != nil {
			return nil, err
		}
		result = append(result, parameter...)
		filter = filter[2+qty:]
	}
	return result, nil
}
