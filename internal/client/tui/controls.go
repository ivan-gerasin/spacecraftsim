package tui

import (
	"fmt"

	"github.com/rivo/tview"
)

// Control represents a UI control for a device
type Control interface {
	// GetID returns the device ID
	GetID() string
	// GetLabel returns the display label
	GetLabel() string
	// GetValue returns the current value
	GetValue() string
	// SetValue sets the current value
	SetValue(value string)
	// GetFormItem returns the tview form item
	GetFormItem() tview.FormItem
}

// BaseControl provides common functionality for controls
type BaseControl struct {
	id    string
	label string
}

func (c *BaseControl) GetID() string    { return c.id }
func (c *BaseControl) GetLabel() string { return c.label }

// CheckboxControl implements a checkbox control
type CheckboxControl struct {
	*BaseControl
	checkbox *tview.Checkbox
	onChange func(checked bool)
}

func NewCheckboxControl(id, label string, onChange func(checked bool)) *CheckboxControl {
	c := &CheckboxControl{
		BaseControl: &BaseControl{id: id, label: label},
		onChange:    onChange,
	}
	c.checkbox = tview.NewCheckbox().SetLabel(label)
	c.checkbox.SetChangedFunc(onChange)
	return c
}

func (c *CheckboxControl) GetValue() string            { return fmt.Sprintf("%v", c.checkbox.IsChecked()) }
func (c *CheckboxControl) SetValue(value string)       { c.checkbox.SetChecked(value == "true") }
func (c *CheckboxControl) GetFormItem() tview.FormItem { return c.checkbox }

// SelectorControl implements a dropdown control
type SelectorControl struct {
	*BaseControl
	dropdown *tview.DropDown
	options  []string
	onChange func(option string, index int)
}

func NewSelectorControl(id, label string, options []string, onChange func(option string, index int)) *SelectorControl {
	c := &SelectorControl{
		BaseControl: &BaseControl{id: id, label: label},
		options:     options,
		onChange:    onChange,
	}
	c.dropdown = tview.NewDropDown().SetLabel(label)
	c.dropdown.SetOptions(options, onChange)
	return c
}

func (c *SelectorControl) GetValue() string {
	_, value := c.dropdown.GetCurrentOption()
	return value
}
func (c *SelectorControl) SetValue(value string) {
	for i, opt := range c.options {
		if opt == value {
			c.dropdown.SetCurrentOption(i)
			break
		}
	}
}
func (c *SelectorControl) GetFormItem() tview.FormItem { return c.dropdown }

// InputControl implements an input field control
type InputControl struct {
	*BaseControl
	input    *tview.InputField
	onChange func(text string)
}

func NewInputControl(id, label string, onChange func(text string)) *InputControl {
	c := &InputControl{
		BaseControl: &BaseControl{id: id, label: label},
		onChange:    onChange,
	}
	c.input = tview.NewInputField().SetLabel(label)
	c.input.SetChangedFunc(onChange)
	return c
}

func (c *InputControl) GetValue() string            { return c.input.GetText() }
func (c *InputControl) SetValue(value string)       { c.input.SetText(value) }
func (c *InputControl) GetFormItem() tview.FormItem { return c.input }
