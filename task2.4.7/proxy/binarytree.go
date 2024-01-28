package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"strings"
	"time"
)

// TreeNode определяет узел AVL дерева
type TreeNode struct {
	Key    int
	Height int
	Left   *TreeNode
	Right  *TreeNode
}

// AVLTree определяет структуру AVL дерева
type AVLTree struct {
	Root *TreeNode
}

// NewNode создаёт новый узел
func NewNode(key int) *TreeNode {
	return &TreeNode{Key: key, Height: 1}
}

// Insert вставляет ключ в AVL дерево
func (t *AVLTree) Insert(key int) {
	t.Root = insert(t.Root, key)
}

// height возвращает высоту узла
func height(node *TreeNode) int {
	if node == nil {
		return 0
	}
	return node.Height
}

// max возвращает максимальное значение из двух чисел
func max(a, b int) int {
	return int(math.Max(float64(a), float64(b)))
}

// updateHeight обновляет высоту узла
func updateHeight(node *TreeNode) {
	node.Height = max(height(node.Left), height(node.Right)) + 1
}

// getBalance возвращает баланс узла
func getBalance(node *TreeNode) int {
	if node == nil {
		return 0
	}
	return height(node.Left) - height(node.Right)
}

// leftRotate выполняет левый поворот
func leftRotate(x *TreeNode) *TreeNode {
	y := x.Right
	T2 := y.Left

	// Выполняем поворот
	y.Left = x
	x.Right = T2

	// Обновляем высоты
	updateHeight(x)
	updateHeight(y)

	// Возвращаем новый корень
	return y
}

// rightRotate выполняет правый поворот
func rightRotate(y *TreeNode) *TreeNode {
	x := y.Left
	T2 := x.Right

	// Выполняем поворот
	x.Right = y
	y.Left = T2

	// Обновляем высоты
	updateHeight(y)
	updateHeight(x)

	// Возвращаем новый корень
	return x
}

// insert вставляет ключ в дерево и выполняет балансировку
func insert(node *TreeNode, key int) *TreeNode {
	// Обычная вставка BST
	if node == nil {
		return NewNode(key)
	}

	if key < node.Key {
		node.Left = insert(node.Left, key)
	} else if key > node.Key {
		node.Right = insert(node.Right, key)
	} else {
		// Дубликаты ключей не допускаются
		return node
	}

	// Обновляем высоту этого узла
	updateHeight(node)

	// Получаем баланс этого узла
	balance := getBalance(node)

	// Если узел несбалансирован, то 4 случая

	// Left Left Case
	if balance > 1 && key < node.Left.Key {
		return rightRotate(node)
	}

	// Right Right Case
	if balance < -1 && key > node.Right.Key {
		return leftRotate(node)
	}

	// Left Right Case
	if balance > 1 && key > node.Left.Key {
		node.Left = leftRotate(node.Left)
		return rightRotate(node)
	}

	// Right Left Case
	if balance < -1 && key < node.Right.Key {
		node.Right = rightRotate(node.Right)
		return leftRotate(node)
	}

	// Возвращаем неизмененный узел
	return node
}

func treeSize(node *TreeNode) int {
	if node == nil {
		return 0
	}
	return 1 + treeSize(node.Left) + treeSize(node.Right)
}

func GenerateTree(count int) *AVLTree {
	tree := &AVLTree{}

	// Инициализируем генератор случайных чисел
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < count; i++ {
		// Генерируем случайный ключ для узла
		key := rand.Intn(100) // Здесь можно настроить диапазон значений
		tree.Insert(key)
	}

	return tree
}
func generateMermaid(node *TreeNode, builder *strings.Builder) {
	if node != nil {
		// Если есть левый потомок, добавляем связь и рекурсивно обрабатываем левого потомка
		if node.Left != nil {
			builder.WriteString(fmt.Sprintf("    %d --> %d\n", node.Key, node.Left.Key))
			generateMermaid(node.Left, builder)
		}
		// Если есть правый потомок, добавляем связь и рекурсивно обрабатываем правого потомка
		if node.Right != nil {
			builder.WriteString(fmt.Sprintf("    %d --> %d\n", node.Key, node.Right.Key))
			generateMermaid(node.Right, builder)
		}
	}
}

// ToMermaid создает строку в формате Mermaid из AVL дерева
func (t *AVLTree) ToMermaid() string {
	var builder strings.Builder
	builder.WriteString("{{< mermaid >}}\ngraph TD\n")
	generateMermaid(t.Root, &builder)
	builder.WriteString("{{< /mermaid >}}\n")
	return builder.String()
}

func updateTree(filePath string, tree *AVLTree) error {
	// Читаем содержимое файла
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	contentStr := string(file)

	// Находим маркеры начала и конца дерева
	startMarker := "<!-- AVLTreeStart -->"
	endMarker := "<!-- AVLTreeEnd -->"
	startIdx := strings.Index(contentStr, startMarker)
	endIdx := strings.Index(contentStr, endMarker)

	if startIdx == -1 || endIdx == -1 {
		return fmt.Errorf("Markers not found in the file")
	}

	// Получаем новое представление дерева
	newTreeContent := tree.ToMermaid()

	// Обновляем содержимое файла
	newContent := contentStr[:startIdx] + startMarker + "\n" + newTreeContent + "\n" + endMarker + contentStr[endIdx+len(endMarker):]

	// Перезаписываем файл с новым содержимым
	err = ioutil.WriteFile(filePath, []byte(newContent), 0644)
	if err != nil {
		return err
	}

	return nil
}
