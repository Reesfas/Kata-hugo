---
menu:
    after:
        name: graph
        weight: 1
title: Построение графа
---

# Построение графа

Нужно написать воркер, который будет строить граф на текущей странице, каждые 5 секунд
От 5 до 30 элементов, случайным образом. Все ноды графа должны быть связаны.

```go
type Node struct {
    ID int
    Name string
	Form string // "circle", "rect", "square", "ellipse", "round-rect", "rhombus"
    Links []*Node
}
```

## Mermaid Chart

[MermaidJS](https://mermaid-js.github.io/) is library for generating svg charts and diagrams from text.

## Пример

{{< columns >}}
```tpl
{{</*/* mermaid [class="text-center"]*/*/>}}
graph LR
A[Square Rect] --> B((Circle))
A --> C(Round Rect)
B --> D{Rhombus}
C --> D
C --> B
{{</*/* /mermaid */*/>}}
```

<--->

{{< mermaid >}}
graph LR
A[Square Rect] --> B((Circle))
A --> C(Round Rect)
B --> D{Rhombus}
C --> D
C --> B
{{< /mermaid >}}

{{< /columns >}}

<!-- GraphStart -->
graph LR
{{< mermaid >}}
graph LR
Node0[Square Rect]
Node0 --> Node4
Node0 --> Node15
Node0 --> Node18
Node1((Circle))
Node1 --> Node3
Node1 --> Node23
Node2[Square Rect]
Node2 --> Node16
Node3[Square Rect]
Node3 --> Node1
Node4[Square Rect]
Node4 --> Node24
Node4 --> Node22
Node5{Rhombus}
Node5 --> Node10
Node6(Round Rect)
Node6 --> Node8
Node7((Circle))
Node7 --> Node21
Node7 --> Node13
Node7 --> Node6
Node8[Square Rect]
Node8 --> Node9
Node8 --> Node9
Node8 --> Node3
Node9((Circle))
Node9 --> Node11
Node10((Circle))
Node10 --> Node15
Node10 --> Node22
Node11[Square Rect]
Node11 --> Node12
Node11 --> Node10
Node12{Rhombus}
Node12 --> Node20
Node13((Circle))
Node13 --> Node10
Node13 --> Node19
Node14[Square Rect]
Node14 --> Node24
Node14 --> Node17
Node14 --> Node1
Node15((Circle))
Node15 --> Node4
Node15 --> Node0
Node16((Circle))
Node16 --> Node10
Node17(Round Rect)
Node17 --> Node0
Node17 --> Node1
Node18[Square Rect]
Node18 --> Node5
Node18 --> Node8
Node18 --> Node0
Node19((Circle))
Node19 --> Node4
Node19 --> Node12
Node19 --> Node4
Node20((Circle))
Node20 --> Node6
Node21((Circle))
Node21 --> Node12
Node22((Circle))
Node22 --> Node17
Node22 --> Node8
Node22 --> Node12
Node23((Circle))
Node23 --> Node14
Node24((Circle))
Node24 --> Node3
Node24 --> Node0
{{< /mermaid >}}
<!-- GraphEnd -->
