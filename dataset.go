package main

import (
	"github.com/scritchley/orc"
	"os"
	//"compress/flate"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
)

func save_ct(filename string, model *gui.QStandardItemModel) {

	filename = getrootitemname(model)

	filename = filename + `.ct`

	f, err := os.Create(filename)
	if err != nil {
		println(err.Error())
	}

	defer f.Close()

	schema, err := orc.ParseSchema(`struct<id:int,termid:string,id_parent:int>`)
	if err != nil {
		println(err.Error())
	}

	w, err := orc.NewWriter(f, orc.SetSchema(schema)) ///w, err := orc.NewWriter(f, orc.SetSchema(schema), orc.SetCompression(orc.CompressionZlib{Level: flate.DefaultCompression}))
	if err != nil {
		println(err.Error())
	}

	root := model.InvisibleRootItem().Index()

	uid := int64(0)
	saveforeach(model, root, w, 0, &uid)

	err = w.Close()
	if err != nil {
		println(err.Error())
	}

}

func saveforeach(model *gui.QStandardItemModel, parent *core.QModelIndex, w *orc.Writer, pid int64, uid *int64) {

	for i := 0; i < model.RowCount(parent); i++ {

		*uid += 1

		idx := model.Index(i, 0, parent)
		id := *uid
		termid := idx.Data(0).ToString() /// core.Qt__DisplayRole
		id_parent := pid

		err := w.Write(id, termid, id_parent)

		println(id, termid, id_parent)

		if err != nil {
			println(err.Error())
		}

		if model.HasChildren(idx) {
			saveforeach(model, idx, w, id, uid)
		}
	}
}

func load_ct(filename string, model *gui.QStandardItemModel) {

	model.Clear()

	r, err := orc.Open(filename)
	if err != nil {
		println(err.Error())
	}

	c := r.Select("id", "termid", "id_parent")

	var pid int64

	parent := model.InvisibleRootItem()

	path := make(map[int64]*gui.QStandardItem)

	path[0] = parent

	for c.Stripes() {
		for c.Next() {
			id := c.Row()[0].(int64)
			termid := c.Row()[1].(string)
			id_parent := c.Row()[2].(int64)

			if pid != id_parent {

				pid = id_parent
				parent = path[pid]

				child := gui.NewQStandardItem2(termid)
				parent.AppendRow2(child)
				path[id] = child

			} else {

				child := gui.NewQStandardItem2(termid)
				parent.AppendRow2(child)
				path[id] = child

			}
		}
	}

}

func getrootitemname(model *gui.QStandardItemModel) string {

	root := model.InvisibleRootItem().Index()
	idx := model.Index(0, 0, root)

	return idx.Data(0).ToString() /// core.Qt__DisplayRole

}
