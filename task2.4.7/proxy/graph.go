package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
)

type Node struct {
	ID    int
	Name  string
	Form  string // "circle", "rect", "square", "ellipse", "round-rect", "rhombus"
	Links []*Node
}

func generateRandomGraph(min, max int) []*Node {
	nodeCount := rand.Intn(max-min+1) + min
	nodes := make([]*Node, nodeCount)

	// Создание узлов
	for i := 0; i < nodeCount; i++ {
		nodes[i] = &Node{
			ID:   i,
			Name: fmt.Sprintf("Node%d", i),
			Form: randomForm(),
		}
	}

	// Создание случайных связей
	for _, node := range nodes {
		// Определяем количество связей для каждого узла (например, от 1 до 3)
		linkCount := rand.Intn(3) + 1
		for i := 0; i < linkCount; i++ {
			// Выбираем случайный узел для создания связи
			targetIndex := rand.Intn(nodeCount)
			if targetIndex != node.ID {
				node.Links = append(node.Links, nodes[targetIndex])
			}
		}
	}

	return nodes
}

func randomForm() string {
	forms := []string{"circle", "square", "round-rect", "rhombus"}
	return forms[rand.Intn(len(forms))]
}

func updateGraphFile(filePath string, nodes []*Node) error {
	// Читаем существующее содержимое файла
	originalContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("Ошибка при чтении файла: %v", err)
	}

	contentStr := string(originalContent)

	// Находим маркеры начала и конца графа
	startMarker := "<!-- GraphStart -->"
	endMarker := "<!-- GraphEnd -->"
	startIdx := strings.Index(contentStr, startMarker)
	endIdx := strings.Index(contentStr, endMarker)

	if startIdx == -1 || endIdx == -1 {
		return fmt.Errorf("Markers not found in the file")
	}

	// Генерируем новое содержимое для графа
	mermaidContent := generateMermaidContent(nodes)

	// Обновляем содержимое файла
	newContent := contentStr[:startIdx] + startMarker + "\ngraph LR\n" + mermaidContent + endMarker + contentStr[endIdx+len(endMarker):]

	// Перезаписываем файл с новым содержимым
	err = ioutil.WriteFile(filePath, []byte(newContent), 0644)
	if err != nil {
		return fmt.Errorf("Ошибка при записи в файл: %v", err)
	}

	return nil
}

func generateMermaidContent(nodes []*Node) string {
	mermaidContent := "{{< mermaid >}}\ngraph LR\n"

	for _, node := range nodes {
		var nodeFormat string
		switch node.Form {
		case "square":
			nodeFormat = "[Square Rect]"
		case "circle":
			nodeFormat = "((Circle))"
		case "round-rect":
			nodeFormat = "(Round Rect)"
		case "rhombus":
			nodeFormat = "{Rhombus}"
		}

		formattedNode := node.Name + nodeFormat
		mermaidContent += formattedNode + "\n"

		for _, link := range node.Links {
			mermaidContent += fmt.Sprintf("%s --> %s\n", node.Name, link.Name)
		}
	}

	mermaidContent += "{{< /mermaid >}}\n"
	return mermaidContent
}
