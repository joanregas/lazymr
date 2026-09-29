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


package gui

import "lazymr/pkg/boxlayout"

type WindowArrangement struct{}

func NewWindowArrangement() *WindowArrangement {
	return &WindowArrangement{}
}

func (w *WindowArrangement) GetWindowDimensions(
	width int,
	height int,
	commandLogsVisible bool,
) map[string]boxlayout.Dimensions {
	rightSideChildren := []*boxlayout.Box{
		{
			Window: overviewView,
			Weight: 8,
		},
	}

	if commandLogsVisible {
		rightSideChildren = append(
			rightSideChildren,
			&boxlayout.Box{
				Window: commandLogsView,
				Weight: 2,
			},
		)
	}

	root := &boxlayout.Box{
		Direction: boxlayout.ROW,

		Children: []*boxlayout.Box{
			{
				Direction: boxlayout.COLUMN,
				Weight:    1,

				Children: []*boxlayout.Box{
					{
						Direction: boxlayout.ROW,
						Weight:    1,

						Children: []*boxlayout.Box{
							{
								Window: repositoriesView,
								Size:   3,
							},
							{
								Window: mergeRequestView,
								Weight: 2,
							},
							{
								Window: filesView,
								Weight: 2,
							},
							{
								Window: pipelinesView,
								Weight: 2,
							},
							{
								Window: statusView,
								Weight: 1,
							},
						},
					},
					{
						Direction: boxlayout.ROW,
						Weight:    2,

						Children: rightSideChildren,
					},
				},
			},
			{
				Window: footerView,
				Size:   1,
			},
		},
	}

	return boxlayout.ArrangeWindows(
		root,
		0,
		0,
		width,
		height,
	)
}


func commandLogsLayout(
	visible bool,
) []*boxlayout.Box {
	if !visible {
		return []*boxlayout.Box{
			{
				Window: overviewView,
				Weight: 1,
			},
		}
	}

	return []*boxlayout.Box{
		{
			Window: overviewView,
			Weight: 8,
		},
		{
			Window: commandLogsView,
			Weight: 2,
		},
	}
}

