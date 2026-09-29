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

import (
	"fmt"

	"github.com/jesseduffield/gocui"
	"lazymr/pkg/boxlayout"
)

func (gui *GUI) setViewDimensions(
	name string,
	dimensions boxlayout.Dimensions,
) error {
	v, err := gui.g.View(name)
	if err != nil {
		return fmt.Errorf(
			"get view %q: %w",
			name,
			err,
		)
	}

	frameOffset := 1
	if v.Frame {
		frameOffset = 0
	}

	_, err = gui.g.SetView(
		name,
		dimensions.X0-frameOffset,
		dimensions.Y0-frameOffset,
		dimensions.X1+frameOffset,
		dimensions.Y1+frameOffset,
		0,
	)

	if err != nil &&
		err != gocui.ErrUnknownView {
		return fmt.Errorf(
			"set dimensions for view %q: %w",
			name,
			err,
		)
	}

	return nil
}

func (gui *GUI) layout(g *gocui.Gui) error {
	width, height := g.Size()

	if width <= 0 || height <= 0 {
		return nil
	}

	dimensions := gui.windowArrangement.GetWindowDimensions(
		width,
		height,
		gui.commandLogsVisible,
	)

	for name, d := range dimensions {
		if err := gui.setViewDimensions(
			name,
			d,
		); err != nil {
			return err
		}
	}

	if !gui.commandLogsVisible {
		if _, err := gui.g.SetView(
			commandLogsView,
			0,
			0,
			0,
			0,
			0,
		); err != nil {
			return fmt.Errorf(
				"hide command logs view: %w",
				err,
			)
		}
	}

	return nil
}


