package main

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

type EditTreeDelegate struct{ widgets.QStyledItemDelegate }

type EditTree struct {
	*widgets.QTreeView
	Im       *gui.QStandardItemModel
	EditMode bool

	editbutton *widgets.QPushButton
}

func New_EditTree(window *widgets.QMainWindow) (*EditTree, *widgets.QVBoxLayout) {

	edittree := &EditTree{QTreeView: widgets.NewQTreeView(window)}
	model := gui.NewQStandardItemModel(window)
	edittree.SetModel(model)
	edittree.Im = model
	edittree.SetFont(font)
	edittree.SetStyleSheet(edittreestylesheet)

	header := gui.NewQStandardItem2(`Classification Tree`)
	header.SetFont(font)
	model.SetHorizontalHeaderItem(0, header)

	edittree.SetDragDropMode(widgets.QAbstractItemView__DragDrop)
	edittree.SetDefaultDropAction(core.Qt__MoveAction)
	edittree.Viewport().SetAcceptDrops(true)
	edittree.SetAcceptDrops(true)
	edittree.SetDropIndicatorShown(true)

	delegate := InitEditTreeDelegate()
	edittree.SetItemDelegate(delegate)

	edittree.SetEditTriggers(widgets.QAbstractItemView__NoEditTriggers)

	hlayout := edittree_ui(edittree, window)

	vlayout := widgets.NewQVBoxLayout()

	vlayout.AddLayout(hlayout, 0)
	vlayout.AddWidget(edittree, 0, 0)

	return edittree, vlayout

}

func edittree_ui(edittree *EditTree, window *widgets.QMainWindow) *widgets.QHBoxLayout {

	hlayout := widgets.NewQHBoxLayout()

	addbutton := widgets.NewQPushButton3(gui.NewQIcon5(":/icons/addnode.png"), "   Add", window)
	addbutton.ConnectClicked(func(_ bool) { edittree.addbutton_click() })

	openbutton := widgets.NewQPushButton3(gui.NewQIcon5(":/icons/opentree.png"), "   Open", window)
	openbutton.ConnectClicked(func(_ bool) { edittree.openbutton_click() })

	clearbutton := widgets.NewQPushButton3(gui.NewQIcon5(":/icons/cleartree.png"), "   Clear", window)
	clearbutton.ConnectClicked(func(_ bool) { edittree.clearbutton_click() })

	savebutton := widgets.NewQPushButton3(gui.NewQIcon5(":/icons/savetree.png"), "   Save", window)
	savebutton.ConnectClicked(func(_ bool) { edittree.savebutton_click() })

	editbutton := widgets.NewQPushButton3(gui.NewQIcon5(":/icons/edittree.png"), " Edit: OFF", window)
	editbutton.ConnectClicked(func(_ bool) { edittree.editbutton_click() })

	edittree.editbutton = editbutton

	hlayout.AddWidget(addbutton, 0, 0)
	hlayout.AddWidget(openbutton, 0, 0)
	hlayout.AddWidget(clearbutton, 0, 0)
	hlayout.AddWidget(savebutton, 0, 0)
	hlayout.AddWidget(editbutton, 0, 0)

	return hlayout
}

func (edittree *EditTree) editbutton_click() {

	if edittree.EditMode {
		edittree.EditMode = false
		edittree.editbutton.SetText(" Edit: OFF")
		edittree.SetEditTriggers(widgets.QAbstractItemView__NoEditTriggers)

	} else {
		edittree.EditMode = true
		edittree.editbutton.SetText(" Edit: ON")
		edittree.SetEditTriggers(widgets.QAbstractItemView__DoubleClicked)
	}

}

func (edittree *EditTree) addbutton_click() {

	edittree.Im.AppendRow2(gui.NewQStandardItem2(`new node`))

}

func (edittree *EditTree) openbutton_click() {

	fd := widgets.NewQFileDialog(ap.Window, core.Qt__Dialog)
	fd.SetViewMode(widgets.QFileDialog__Detail)

	fn := fd.GetOpenFileName(ap.Window, `Open File`, core.QDir_CurrentPath(), `Classify Trees (*.ct)`, "", widgets.QFileDialog__ReadOnly)
	if fn == `` {
		return
	}

	model := edittree.Im
	load_ct(fn, model)

	header := gui.NewQStandardItem2(getrootitemname(model))
	header.SetFont(font)
	model.SetHorizontalHeaderItem(0, header)

}

func (edittree *EditTree) clearbutton_click() {
	edittree.Im.Clear()
}

func (edittree *EditTree) savebutton_click() {

	model := edittree.Im
	save_ct("", model)

}

func InitEditTreeDelegate() *EditTreeDelegate {
	item := NewEditTreeDelegate(nil)
	item.ConnectCreateEditor(createEditor_et)
	item.ConnectSetEditorData(setEditorData_et)
	item.ConnectSetModelData(setModelData_et)
	item.ConnectUpdateEditorGeometry(updateEditorGeometry_et)
	return item
}

func createEditor_et(parent *widgets.QWidget, option *widgets.QStyleOptionViewItem, index *core.QModelIndex) *widgets.QWidget {

	editor := widgets.NewQLineEdit(parent)

	return editor.QWidget_PTR()
}

func setEditorData_et(editor *widgets.QWidget, index *core.QModelIndex) {

	lineedit := widgets.NewQLineEditFromPointer(editor.Pointer())

	value := index.Model().Data(index, int(core.Qt__EditRole)).ToString()
	lineedit.SetText(value)
}

func setModelData_et(editor *widgets.QWidget, model *core.QAbstractItemModel, index *core.QModelIndex) {

	lineedit := widgets.NewQLineEditFromPointer(editor.Pointer())

	text := lineedit.Text()
	model.SetData(index, core.NewQVariant14(text), int(core.Qt__EditRole))

}

func updateEditorGeometry_et(editor *widgets.QWidget, option *widgets.QStyleOptionViewItem, index *core.QModelIndex) {
	editor.SetGeometry(option.Rect())
}

var edittreestylesheet string = `
	QTreeView::branch:has-siblings:!adjoins-item {
		 border-image: url(:/tree/vline.png) 0;
	 }

	 QTreeView::branch:has-siblings:adjoins-item {
		 border-image: url(:/tree/branch-more.png) 0;
	 }

	 QTreeView::branch:!has-children:!has-siblings:adjoins-item {
		 border-image: url(:/tree/branch-end.png) 0;
	 }

	 QTreeView::branch:has-children:!has-siblings:closed,
	 QTreeView::branch:closed:has-children:has-siblings {
			 border-image: none;
			 image: url(:/tree/branch-closed.png);
	 }

	 QTreeView::branch:open:has-children:!has-siblings,
	 QTreeView::branch:open:has-children:has-siblings  {
			 border-image: none;
			 image: url(:/tree/branch-open.png);
	 }
	 

	QTreeView {
			selection-background-color: lightGrey;
	}

	QTreeView::item:selected
	{		background-color: lightGrey;
			color: black;		
	}
	 
	 
`
