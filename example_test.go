package main

import (
	"context"
	"log"
	"testing"
	"time"

	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/model"
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/tcvectordb"
)

var cli model.VectorDBClient

func init() {
	var err error
	cli, err = tcvectordb.NewClient("http://11.141.218.159:8100", "root", "p193scgeHBYRlDWHKCVfIm5z8eI0Q96HyArkbqNg", &model.ClientOption{
		MaxIdldConnPerHost: 50,
		IdleConnTimeout:    time.Second * 30,
	})
	if err != nil {
		panic(err)
	}
}

func TestDatabase(t *testing.T) {
	defer cli.Close()

	_, err := cli.CreateDatabase(context.TODO(), "dbtest1")
	printErr(err)

	_, err = cli.CreateDatabase(context.TODO(), "dbtest2")
	printErr(err)

	dbs, err := cli.ListDatabase(context.Background())
	printErr(err)

	for _, db := range dbs {
		t.Logf("database: %s", db.DatabaseName)
	}
	err = cli.DropDatabase(context.Background(), "dbtest2")
	printErr(err)

	dbs, err = cli.ListDatabase(context.Background())
	printErr(err)
	for _, db := range dbs {
		t.Logf("database: %s", db.DatabaseName)
	}
}

func TestCreateCollection(t *testing.T) {
	defer cli.Close()

	db := cli.Database("dbtest1")
	_ = db.DropCollection(context.Background(), "col1")

	_, err := db.CreateCollection(context.Background(), "col1", 2, 2, "desription doc", model.Indexes{
		VectorIndex: []model.VectorIndex{
			{
				FilterIndex: model.FilterIndex{
					FieldName: "vector",
					FieldType: model.Vector,
					IndexType: model.HNSW,
				},
				Dimension:  3,
				MetricType: model.L2,
				HNSWParam: model.HNSWParam{
					M:              64,
					EfConstruction: 8,
				},
			},
		},
		FilterIndex: []model.FilterIndex{
			{
				FieldName: "id",
				FieldType: model.String,
				IndexType: model.PRIMARY,
			},
			{
				FieldName: "author",
				FieldType: model.String,
				IndexType: model.FILTER,
			},
			{
				FieldName: "page",
				FieldType: model.Uint64,
				IndexType: model.FILTER,
			},
		},
	}, nil)
	printErr(err)

	colList, err := db.ListCollection(context.TODO())
	printErr(err)
	for _, col := range colList {
		t.Logf("%+v", col)
	}

	col, err := db.DescribeCollection(context.Background(), "col1")
	printErr(err)
	t.Logf("%+v", col)
}

func TestCreateCollectionWithEmbedding(t *testing.T) {
	defer cli.Close()

	db := cli.Database("dbtest1")
	_ = db.DropCollection(context.Background(), "col2")

	em := model.Embedding{
		TextField:   "text",
		VectorField: "vector",
		Model:       model.M3E_BASE,
		Enabled:     true,
	}

	cli.Debug(true)
	_, err := db.CreateCollection(context.Background(), "col2", 2, 2, "desription doc", model.Indexes{
		VectorIndex: []model.VectorIndex{
			{
				FilterIndex: model.FilterIndex{
					FieldName: "vector",
					FieldType: model.Vector,
					IndexType: model.HNSW,
				},
				Dimension:  768,
				MetricType: model.L2,
				HNSWParam: model.HNSWParam{
					M:              64,
					EfConstruction: 8,
				},
			},
		},
		FilterIndex: []model.FilterIndex{
			{
				FieldName: "id",
				FieldType: model.String,
				IndexType: model.PRIMARY,
			},
			{
				FieldName: "author",
				FieldType: model.String,
				IndexType: model.FILTER,
			},
			{
				FieldName: "page",
				FieldType: model.Uint64,
				IndexType: model.FILTER,
			},
		},
	}, &em)
	printErr(err)

	col, err := db.DescribeCollection(context.Background(), "col2")
	printErr(err)
	t.Logf("%+v", col)
}

func TestGetCollection(t *testing.T) {
	defer cli.Close()

	db := cli.Database("dbtest1")
	cli.Debug(true)
	col, err := db.DescribeCollection(context.TODO(), "col2")
	printErr(err)
	t.Logf("%+v", col)
}

func TestAlias(t *testing.T) {
	db := cli.Database("dbtest1")
	// db.Debug(true)
	var affectCount int
	var err error
	affectCount, err = db.AliasSet(context.Background(), "col1", "alias-col1")
	t.Logf("affect count: %d", affectCount)
	printErr(err)
	affectCount, err = db.AliasSet(context.Background(), "col2", "alias-col2")
	t.Logf("affect count: %d", affectCount)
	printErr(err)

	// aliasList, err := db.AliasList(context.Background())
	// printErr(err)
	// for _, item := range aliasList {
	// 	t.Logf("%+v", item)
	// }

	// alias, err := db.AliasDescribe(context.Background(), "alias-col1")
	// printErr(err)
	// t.Logf("%+v", alias)
	// affectCount, err = db.AliasDrop(context.Background(), "alias-col1")
	// t.Logf("affect count: %d", affectCount)
	// printErr(err)
	// affectCount, err = db.AliasDrop(context.Background(), "alias-col2")
	// t.Logf("affect count: %d", affectCount)
	// printErr(err)
	// aliasList, err = db.AliasList(context.Background())
	// printErr(err)
	// for _, item := range aliasList {
	// 	t.Logf("%+v", item)
	// }
}

func TestIndex(t *testing.T) {
	db := cli.Database("dbtest1")
	err := db.IndexRebuild(context.Background(), "col1", false, 1)
	printErr(err)
}

func TestUpsertDocument(t *testing.T) {
	defer cli.Close()
	col := cli.Database("dbtest1").Collection("col1")

	err := col.Upsert(context.Background(), []model.Document{
		{
			Id:     "0001",
			Vector: []float32{0.2123, 0.23, 0.213},
			Fields: map[string]model.Field{
				"author":  model.Field{Val: "jerry"},
				"page":    model.Field{Val: 21},
				"section": model.Field{Val: "1.1.1"},
			},
		},
		{
			Id:     "0002",
			Vector: []float32{0.2123, 0.22, 0.213},
			Fields: map[string]model.Field{
				"author":  model.Field{Val: "sam"},
				"page":    model.Field{Val: 22},
				"section": model.Field{Val: "1.1.2"},
			},
		},
		{
			Id:     "0003",
			Vector: []float32{0.2123, 0.21, 0.213},
			Fields: map[string]model.Field{
				"author":  model.Field{Val: "max"},
				"page":    model.Field{Val: 23},
				"section": model.Field{Val: "1.1.3"},
			},
		},
	}, true)

	printErr(err)
}

func TestSearch(t *testing.T) {
	defer cli.Close()
	col := cli.Database("dbtest1").Collection("col1")
	t.Log("document query-----------------")
	docs, count, err := col.Query(context.Background(), []string{"0001", "0002"}, nil, "", true, nil, 0, 10)
	printErr(err)
	t.Logf("total doc: %d", count)
	for _, doc := range docs {
		t.Logf("id: %s, vector: %v, author: %s, page: %d, section: %s", doc.Id, doc.Vector,
			doc.Fields["author"].String(), doc.Fields["page"].Int(), doc.Fields["section"].String())
	}
	t.Log("document search-----------------")
	filter := model.NewFilter("page > 22").And(model.In("author", []string{"max", "sam"}))
	searchRes, err := col.Search(context.Background(), [][]float32{{0.3123, 0.43, 0.213}}, nil, filter, "", &model.HNSWParam{EfConstruction: 10}, true, nil, 10)
	printErr(err)
	for i, docs := range searchRes {
		t.Logf("doc %d result: ", i)
		for _, doc := range docs {
			t.Logf("id: %s, vector: %v, score: %v, author: %s, page: %d, section: %s", doc.Id, doc.Vector, doc.Score,
				doc.Fields["author"].String(), doc.Fields["page"].Int(), doc.Fields["section"].String())
		}
	}

	col.Debug(true)
	t.Log("document searchById-----------------")
	searchRes, err = col.SearchById(context.Background(), []string{"0001", "0002", "0003"}, nil, filter, "", &model.HNSWParam{EfConstruction: 10}, true, nil, 10)
	printErr(err)
	for i, docs := range searchRes {
		t.Logf("doc %d result: ", i)
		for _, doc := range docs {
			t.Logf("id: %s, vector: %v, author: %s, page: %d, section: %s", doc.Id, doc.Vector,
				doc.Fields["author"].String(), doc.Fields["page"].Int(), doc.Fields["section"].String())
		}
	}
}

func TestDeleteDocument(t *testing.T) {
	var err error
	defer cli.Close()
	db := cli.Database("dbtest1")
	col := db.Collection("col1")

	// delete documents
	err = col.Delete(context.Background(), []string{"0002", "0003"}, nil)
	printErr(err)
	err = col.Delete(context.Background(), []string{"0002", "0003"}, nil)
	printErr(err)

	docs, count, err := col.Query(context.Background(), []string{"0002", "0003"}, nil, "", false, nil, 0, 10)
	t.Logf("affect doc: %d", count)
	if len(docs) != 0 {
		t.Errorf("%v", docs)
	}

}

func TestDeleteCollection(t *testing.T) {
	var err error
	defer cli.Close()
	db := cli.Database("dbtest1")
	err = db.DropCollection(context.Background(), "col1")
	printErr(err)
	err = db.DropCollection(context.Background(), "col1")
	printErr(err)
	collist, err := db.ListCollection(context.Background())
	printErr(err)
	for _, col := range collist {
		t.Logf("%v", col)
	}
}

func TestDeleteDatabase(t *testing.T) {
	var err error
	defer cli.Close()
	err = cli.DropDatabase(context.Background(), "dbtest1")
	printErr(err)
	err = cli.DropDatabase(context.Background(), "dbtest1")
	printErr(err)
	dbs, err := cli.ListDatabase(context.Background())
	printErr(err)
	for _, db := range dbs {
		t.Log(db.DatabaseName)
	}
}

func printErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}