package component

import (

)

type FormTemplateIterator struct {}

type IterateFormTemplateColumnFunc func(column Column, result *interface{})

func (o FormTemplateIterator) IterateTemplateColumn(formTemplate FormTemplate, result *interface{}, iterateFunc IterateFormTemplateColumnFunc) {
	if formTemplate.FormElemLi != nil {
		listTemplateIterator := ListTemplateIterator{}
		for _, item := range formTemplate.FormElemLi {
			if item.XMLName.Local == "column-model" {
				columnLi := []Column{}
				listTemplateIterator.recursionGetColumnItem(item.ColumnModel, &columnLi)
				for _, item := range columnLi {
					iterateFunc(item, result)
				}
			}
		}
	}
}

type IterateTemplateColumnModelFunc func(columnModel ColumnModel, result * interface{})

func (o FormTemplateIterator) IterateAllTemplateColumnModel(formTemplate FormTemplate, result *interface{}, iterateFunc IterateTemplateColumnModelFunc) {
	for _, item := range formTemplate.FormElemLi {
		if (item.XMLName.Local == "column-model") {
			iterateFunc(item.ColumnModel, result);
		}
	}
}

type IterateTemplateButtonFunc func(toolbar Toolbar, column ColumnModel, isToolbarBtn bool, button Button, result *interface{}) 

/**
 * 按钮的分布,toolbar/button,columnModel/toolbar,columnModel/editor-toolbar,columnModel/virtual-column/buttons/button
 */
func (o FormTemplateIterator) IterateTemplateButton(formTemplate FormTemplate, result *interface{}, iterateFunc IterateTemplateButtonFunc) {
	for j, _ := range formTemplate.FormElemLi {
		formElem := formTemplate.FormElemLi[j]
		if formElem.XMLName.Local == "toolbar" {
			if formElem.Toolbar.ButtonLi != nil {
				for k, _ := range formElem.Toolbar.ButtonLi {
					button := formElem.Toolbar.ButtonLi[k]
					isToolbarBtn := true
					iterateFunc(formElem.Toolbar, ColumnModel{}, isToolbarBtn, button, result)
				}
			}
		} else if formElem.XMLName.Local == "column-model" {
			if formElem.ColumnModel.Toolbar.ButtonLi != nil {
				for k, _ := range formElem.ColumnModel.Toolbar.ButtonLi {
					button := formElem.ColumnModel.Toolbar.ButtonLi[k]
					isToolbarBtn := false
					iterateFunc(Toolbar{}, formElem.ColumnModel, isToolbarBtn, button, result)
				}
			}
			if formElem.ColumnModel.EditorToolbar.ButtonLi != nil {
				for k, _ := range formElem.ColumnModel.EditorToolbar.ButtonLi {
					button := formElem.ColumnModel.EditorToolbar.ButtonLi[k]
					isToolbarBtn := false
					iterateFunc(Toolbar{}, formElem.ColumnModel, isToolbarBtn, button, result)
				}
			}
			if formElem.ColumnModel.ColumnLi != nil {
				for k, _ := range formElem.ColumnModel.ColumnLi {
					column := formElem.ColumnModel.ColumnLi[k]
					if column.XMLName.Local == "virtual-column" {
						if column.Buttons.ButtonLi != nil {
							for l, _ := range column.Buttons.ButtonLi {
								button := column.Buttons.ButtonLi[l]
								isToolbarBtn := false
								iterateFunc(Toolbar{}, formElem.ColumnModel, isToolbarBtn, button, result)
							}
						}
					}
				}
			}
		}
	}
}

