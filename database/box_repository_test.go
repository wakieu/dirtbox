package database

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/wakieu/drtbox/entity"
)

func setupTest(tb testing.TB) (*BoxRepository, func(tb testing.TB)) {
	db, err := sql.Open("sqlite3", "./test_db.sqlite")
	if err != nil {
		tb.Error(err)
	}

	SQL := "CREATE TABLE box (boxpath TEXT NOT NULL PRIMARY KEY, text TEXT);"
	_, err = db.Exec(SQL)
	if err != nil && err.Error() != "table box already exists" {
		log.Printf("%q\n", err)
	}

	teardownTest := func(tb testing.TB) {
		err := os.Remove("./test_db.sqlite")
		if err != nil {
			log.Printf("Test teardown failed! (%v)", err)
		}
		defer db.Close()
	}

	return NewBoxRepository(db), teardownTest
}

func TestBoxRepository_CRUD(t *testing.T) {
	const BOX_PATH = "/foo"
	const BOX_TEXT = "teste 123"

	boxRepo, teardown := setupTest(t)
	defer teardown(t)

	//Test box_repository_Save
	err := boxRepo.Save(entity.NewBoxWithText(BOX_PATH, BOX_TEXT))
	if err != nil {
		t.Fatalf("Error saving box!")
	}

	//Test box_repository_Exists
	ok, err := boxRepo.Exists(BOX_PATH)
	if err != nil {
		t.Fatalf("Error checking existing box!")
	}
	if !ok {
		t.Fatalf("Expected box to exist, got exists false")
	}

	//Test box_repository_GetContent
	box, err := boxRepo.GetContent(BOX_PATH)
	if err != nil {
		t.Fatalf("Error getting box!")
	}
	if box.Text != BOX_TEXT {
		t.Fatalf("Different texts. Expected %v, got %v", BOX_TEXT, box.Text)
	}

	//Test box_repository_GetContent
	err = boxRepo.Delete(BOX_PATH)
	if err != nil {
		t.Fatalf("Error deleting box!")
	}
	ok, err = boxRepo.Exists(BOX_PATH)
	if err != nil {
		t.Fatalf("Error checking deleted box!")
	}
	if ok {
		t.Fatalf("Expected box to have been deleted, got exists true")
	}
}


func TestBoxRepository_Save(t *testing.T) {
	const BOX_PATH = "/foo"

	boxRepo, teardown := setupTest(t)
	defer teardown(t)

	err := boxRepo.Save(entity.NewBox(BOX_PATH))
	if err != nil {
		t.Fatalf("Error saving box!")
	}
	exists, err := boxRepo.Exists(BOX_PATH)
	if err != nil {
		t.Fatalf("Error checking existing box!")
	}
	if exists {
		t.Errorf("Fail! Box with empty text should not be saved!")
	}
}

func TestBoxRepository_Children(t *testing.T) {
	const BOX_PATH_GRANDPARENT = "/foo"
	const BOX_PATH_PARENT = BOX_PATH_GRANDPARENT + "/bar"
	const BOX_PATH_CHILD = BOX_PATH_PARENT + "/baz"
	const BOX_TEXT = "foo bar baz"

	log.Println(BOX_PATH_GRANDPARENT, BOX_PATH_PARENT, BOX_PATH_CHILD)

	boxRepo, teardown := setupTest(t)
	defer teardown(t)

	err := boxRepo.Save(entity.NewBoxWithText(BOX_PATH_CHILD, BOX_TEXT))
	if err != nil {
		t.Fatalf("Error saving box!")
	}
	err = boxRepo.Save(entity.NewBoxWithText(BOX_PATH_PARENT, BOX_TEXT))
	if err != nil {
		t.Fatalf("Error saving box!")
	}
	err = boxRepo.Save(entity.NewBoxWithText(BOX_PATH_GRANDPARENT, BOX_TEXT))
	if err != nil {
		t.Fatalf("Error saving box!")
	}

	grandparentChildren, err := boxRepo.GetChildren(BOX_PATH_GRANDPARENT)
	if err != nil {
		t.Fatalf("Error getting children!")
	}
	if len(grandparentChildren) != 2 {
		t.Errorf("Wrong number of children! Expected %v, got %v. %v", 2, len(grandparentChildren), grandparentChildren)
	}

	parentChildren, err := boxRepo.GetChildren(BOX_PATH_PARENT)
	if err != nil {
		t.Fatalf("Error getting children!")
	}
	if len(parentChildren) != 1 {
		t.Errorf("Wrong number of children! Expected %v, got %v. %v", 1, len(parentChildren),parentChildren)
	}

	childChildren, err := boxRepo.GetChildren(BOX_PATH_CHILD)
	if err != nil {
		t.Fatalf("Error getting children!")
	}
	if len(childChildren) != 0 {
		t.Errorf("Wrong number of children! Expected %v, got %v. %v", 0, len(childChildren), childChildren)
	}
}
