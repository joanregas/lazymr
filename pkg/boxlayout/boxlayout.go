/*
MIT License                                                                                   
                                                                                              
Copyright (c) 2018 Jesse Duffield                                                             
                                                                                              
Permission is hereby granted, free of charge, to any person obtaining a copy                  
of this software and associated documentation files (the "Software"), to deal                 
in the Software without restriction, including without limitation the rights                  
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell                     
copies of the Software, and to permit persons to whom the Software is                         
furnished to do so, subject to the following conditions:                                      
                                                                                              
The above copyright notice and this permission notice shall be included in all                
copies or substantial portions of the Software.                                               
                                                                                              
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR                    
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,                      
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE                   
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER                        
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,                 
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE                 
SOFTWARE.                                                                                     
*/


package boxlayout

type Dimensions struct {
	X0 int
	X1 int
	Y0 int
	Y1 int
}

type Direction int

const (
	ROW Direction = iota
	COLUMN
)

type Box struct {
	Direction            Direction
	ConditionalDirection func(width, height int) Direction

	Children            []*Box
	ConditionalChildren func(width, height int) []*Box

	Window string

	// Size is a fixed size:
	//   ROW    -> height
	//   COLUMN -> width
	Size int

	// Weight is used to distribute remaining space.
	Weight int
}

func ArrangeWindows(
	root *Box,
	x0, y0, width, height int,
) map[string]Dimensions {
	children := root.getChildren(width, height)

	if len(children) == 0 {
		if root.Window == "" {
			return map[string]Dimensions{}
		}

		return map[string]Dimensions{
			root.Window: {
				X0: x0,
				Y0: y0,
				X1: x0 + width - 1,
				Y1: y0 + height - 1,
			},
		}
	}

	direction := root.getDirection(width, height)

	availableSize := height
	if direction == COLUMN {
		availableSize = width
	}

	sizes := calcSizes(children, availableSize)

	result := make(map[string]Dimensions)
	offset := 0

	for i, child := range children {
		boxSize := sizes[i]

		var childResult map[string]Dimensions

		if direction == COLUMN {
			childResult = ArrangeWindows(
				child,
				x0+offset,
				y0,
				boxSize,
				height,
			)
		} else {
			childResult = ArrangeWindows(
				child,
				x0,
				y0+offset,
				width,
				boxSize,
			)
		}

		for name, dimensions := range childResult {
			result[name] = dimensions
		}

		offset += boxSize
	}

	return result
}

func calcSizes(boxes []*Box, availableSpace int) []int {
	weights := make([]int, len(boxes))

	totalWeight := 0
	reservedSpace := 0

	for i, box := range boxes {
		if box.isStatic() {
			reservedSpace += box.Size
			continue
		}

		weight := box.Weight
		if weight < 0 {
			weight = 0
		}

		weights[i] = weight
		totalWeight += weight
	}

	dynamicSpace := availableSpace - reservedSpace
	if dynamicSpace < 0 {
		dynamicSpace = 0
	}

	result := make([]int, len(boxes))

	if totalWeight == 0 {
		for i, box := range boxes {
			if box.isStatic() {
				result[i] = min(availableSpace, box.Size)
			}
		}
		return result
	}

	unitSize := dynamicSpace / totalWeight
	remainder := dynamicSpace % totalWeight

	for i, box := range boxes {
		if box.isStatic() {
			result[i] = min(availableSpace, box.Size)
			continue
		}

		result[i] = unitSize * weights[i]
	}

	// Distribute rounding remainder one character at a time.
	for remainder > 0 {
		for i := range boxes {
			if weights[i] <= 0 {
				continue
			}

			result[i]++
			remainder--

			if remainder == 0 {
				break
			}
		}
	}

	return result
}

func (b *Box) isStatic() bool {
	return b.Size > 0
}

func (b *Box) getDirection(width, height int) Direction {
	if b.ConditionalDirection != nil {
		return b.ConditionalDirection(width, height)
	}

	return b.Direction
}

func (b *Box) getChildren(width, height int) []*Box {
	if b.ConditionalChildren != nil {
		return b.ConditionalChildren(width, height)
	}

	return b.Children
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
