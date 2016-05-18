package api

import (
    "fmt"
    "errors"
    "reflect"
    
    "gopkg.in/olivere/elastic.v3"
	"github.com/Sirupsen/logrus"    

	"github.com/cafebazaar/booker-reservation/common"
)

const (
    esIndex = "reservation"
)

var (
    _elasticClient *elastic.Client
)

func elasticClient() (*elastic.Client, error) {
    if _elasticClient == nil {
        elasticURL := common.ConfigString("ELASTIC_URL")
        if elasticURL == "" {
            return nil, errors.New("No ELASTIC_URL was given")
        }

        client, err := elastic.NewClient(
            elastic.SetSniff(false),
            elastic.SetURL(elasticURL),
            elastic.SetInfoLog(common.LogrusInfoLogger),
            elastic.SetErrorLog(common.LogrusErrorLogger),
        )
        if err != nil {
            return nil, fmt.Errorf("Error while elastic.NewClient: %s", err)
        }

        _elasticClient = client
    }

    return _elasticClient, nil
}

type reservation struct {
    StartTimestamp uint64
    EndTimestamp uint64
    UserID string
    ObjectURI string
} 

func getReservation(objectURI string, timestamp uint64) (*reservation, error) {
    client, err := elasticClient()
    if err != nil {
        return nil, fmt.Errorf("Error while getting elastic client: %s", err)
    }

    exists, err := client.IndexExists(esIndex).Do()
    if err != nil {
        return nil, err
    }
    if !exists {
        return nil, nil
    }

    query := elastic.NewBoolQuery().Must(
        elastic.NewTermQuery("URI", objectURI),
        elastic.NewRangeQuery("startTimestamp").Lt(timestamp),
        elastic.NewRangeQuery("endTimestamp").Gt(timestamp),
    )
    searchResult, err := client.Search().
        Index(esIndex).
        Query(query).
        From(0).Size(1).
        Do()
    if err != nil {
        return nil, fmt.Errorf("Error while querying elastic for reservation at the given time: %s", err)
    }

    var rsv reservation
    for _, item := range searchResult.Each(reflect.TypeOf(rsv)) {
        rsv, ok := item.(reservation)
        if !ok {
            return nil, fmt.Errorf("Failed to convert item to reservation. item=%v", item)
        }

        return &rsv, nil
    }

    return nil, nil
}

func createReservation(objectURI string, startTimestamp, endTimestamp uint64, userID string) (*reservation, error) {
    client, err := elasticClient()
    if err != nil {
        return nil, fmt.Errorf("Error while getting elastic client: %s", err)
    }

    exists, err := client.IndexExists(esIndex).Do()
    if err != nil {
        return nil, fmt.Errorf("Error while checking elastic index: %s", err)
    }
    if !exists {
        createIndex, err := client.CreateIndex(esIndex).Do()
        if err != nil {
            return nil, fmt.Errorf("Error while creating elastic index: %s", err)
        }
        if !createIndex.Acknowledged {
            return nil, errors.New("Error while creating elastic index: Not Acknowledged")
        }
    }

    query := elastic.NewBoolQuery().Must(
        elastic.NewTermQuery("URI", objectURI),
        elastic.NewRangeQuery("startTimestamp").Lt(endTimestamp),
        elastic.NewRangeQuery("endTimestamp").Gt(startTimestamp),
    )
    searchResult, err := client.Search().
        Index(esIndex).
        Query(query).
        From(0).Size(1).
        Pretty(true).
        Do()
    if err != nil {
        return nil, fmt.Errorf("Error while querying elastic for collision: %s", err)
    }

    var rsv reservation
    for _, item := range searchResult.Each(reflect.TypeOf(rsv)) {
        rsv, ok := item.(reservation)
        if !ok {
            logrus.WithField("item", item).Debug("Failed to convert item to reservation")
        }

        return &rsv, errors.New("Collision(s) were found")
    }

    rsv = reservation{
        StartTimestamp: startTimestamp,
        EndTimestamp: endTimestamp,
        UserID: userID,
        ObjectURI: objectURI,
    }

    _, err = client.Index().
        Index(esIndex).
        Type("reservation").
        BodyJson(&rsv).
        Do()
    if err != nil {
        return nil, fmt.Errorf("Error while creating the reservation: %s", err)
    }

    return &rsv, nil
}
