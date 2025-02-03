package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/shopspring/decimal"
	"github.com/vanspaul/SmartMeterSystem/models"
	"github.com/vanspaul/SmartMeterSystem/models/users"

	"github.com/vanspaul/SmartMeterSystem/controllers"

	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	// MongoDB connection details
	uri := "mongodb+srv://vanspaul09:ab7vSvvo14nx7gN3@cluster0.euhiz.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
	dbName := "test_db"
	collName := "persons"

	// Create a new MongoDB controller
	controller, err := controllers.NewMongoDBController(uri, dbName, collName)
	if err != nil {
		log.Fatalf("Failed to create MongoDB controller: %v", err)
	}
	defer func() {
		if err := controller.Close(context.Background()); err != nil {
			log.Printf("Failed to close MongoDB connection: %v", err)
		}
	}()

	// Initialize the controller (optional if using NewMongoDBController)
	if err := controller.Init(context.Background()); err != nil {
		log.Fatalf("Failed to initialize MongoDB controller: %v", err)
	}

	acc := models.Account{
		HashedPassword: "wpefojpanfasdfuivabnib",
		Email:          "johndoe@mail.com",
		CreatedAt:      time.Now().Unix(),
		UpdatedAt:      time.Now().Unix(),
		Role:           models.RoleConsumer,
		RoleSpecificData: users.Consumer{
			AccountNumber: "1234567890",
			Name:          "John Doe",
			Address: users.Address{
				Street:     "Calzada",
				City:       "Balayan",
				State:      "Batangas",
				PostalCode: "4230",
			},
			Contact: users.Contact{
				Phone: "09000000000",
				Email: "johndoe@mail.com",
			},
			Meters: []users.Meter{
				{
					MeterID: "adfbaludifbaufd",
					Location: users.GeoJSON{
						Type:        "point",
						Coordinates: []float64{12.34, 56.78},
					},
					SIM: users.SIM{
						ICCID:          "1830942872875982",
						MobileNumber:   "09000000000",
						DataUsageMb:    506985.0,
						ActivationDate: time.Now().Unix(),
					},
					Usage: []users.Usage{
						{
							Date: 1738490036,
							Kwh:  5.2,
						},
						{
							Date: 1738493636,
							Kwh:  3.8,
						},
						{
							Date: 1738497236,
							Kwh:  7.1,
						},
						{
							Date: 1738500836,
							Kwh:  6.5,
						},
						{
							Date: 1738504436,
							Kwh:  4.9,
						},
						{
							Date: 1738508036,
							Kwh:  2.3,
						},
						{
							Date: 1738511636,
							Kwh:  8.7,
						},
						{
							Date: 1738515236,
							Kwh:  9.4,
						},
						{
							Date: 1738518836,
							Kwh:  1.6,
						},
						{
							Date: 1738522436,
							Kwh:  4.2,
						},
						{
							Date: 1738526036,
							Kwh:  6.8,
						},
						{
							Date: 1738529636,
							Kwh:  5.3,
						},
						{
							Date: 1738533236,
							Kwh:  7.9,
						},
						{
							Date: 1738536836,
							Kwh:  3.1,
						},
						{
							Date: 1738540436,
							Kwh:  2.5,
						},
						{
							Date: 1738544036,
							Kwh:  8.4,
						},
						{
							Date: 1738547636,
							Kwh:  6.7,
						},
						{
							Date: 1738551236,
							Kwh:  4.6,
						},
						{
							Date: 1738554836,
							Kwh:  9.2,
						},
						{
							Date: 1738558436,
							Kwh:  1.8,
						},
						{
							Date: 1738562036,
							Kwh:  5.5,
						},
						{
							Date: 1738565636,
							Kwh:  7.3,
						},
						{
							Date: 1738569236,
							Kwh:  3.9,
						},
						{
							Date: 1738572836,
							Kwh:  2.7,
						},
						{
							Date: 1738576436,
							Kwh:  8.1,
						},
						{
							Date: 1738580036,
							Kwh:  6.4,
						},
						{
							Date: 1738583636,
							Kwh:  4.8,
						},
						{
							Date: 1738587236,
							Kwh:  9.6,
						},
						{
							Date: 1738590836,
							Kwh:  1.3,
						},
						{
							Date: 1738594436,
							Kwh:  5.7,
						},
						{
							Date: 1738598036,
							Kwh:  7.2,
						},
						{
							Date: 1738601636,
							Kwh:  3.4,
						},
						{
							Date: 1738605236,
							Kwh:  2.9,
						},
						{
							Date: 1738608836,
							Kwh:  8.5,
						},
						{
							Date: 1738612436,
							Kwh:  6.1,
						},
						{
							Date: 1738616036,
							Kwh:  4.3,
						},
						{
							Date: 1738619636,
							Kwh:  9.8,
						},
						{
							Date: 1738623236,
							Kwh:  1.5,
						},
						{
							Date: 1738626836,
							Kwh:  5.9,
						},
						{
							Date: 1738630436,
							Kwh:  7.4,
						},
						{
							Date: 1738634036,
							Kwh:  3.6,
						},
						{
							Date: 1738637636,
							Kwh:  2.2,
						},
						{
							Date: 1738641236,
							Kwh:  8.8,
						},
						{
							Date: 1738644836,
							Kwh:  6.3,
						},
						{
							Date: 1738648436,
							Kwh:  4.7,
						},
						{
							Date: 1738652036,
							Kwh:  9.1,
						},
						{
							Date: 1738655636,
							Kwh:  1.9,
						},
						{
							Date: 1738659236,
							Kwh:  5.4,
						},
						{
							Date: 1738662836,
							Kwh:  7.7,
						},
						{
							Date: 1738666436,
							Kwh:  3.2,
						},
						{
							Date: 1738670036,
							Kwh:  2.6,
						},
						{
							Date: 1738673636,
							Kwh:  8.2,
						},
						{
							Date: 1738677236,
							Kwh:  6.6,
						},
						{
							Date: 1738680836,
							Kwh:  4.4,
						},
						{
							Date: 1738684436,
							Kwh:  9.3,
						},
						{
							Date: 1738688036,
							Kwh:  1.4,
						},
						{
							Date: 1738691636,
							Kwh:  5.8,
						},
						{
							Date: 1738695236,
							Kwh:  7.5,
						},
						{
							Date: 1738698836,
							Kwh:  3.7,
						},
						{
							Date: 1738702436,
							Kwh:  2.8,
						},
						{
							Date: 1738706036,
							Kwh:  8.6,
						},
						{
							Date: 1738709636,
							Kwh:  6.2,
						},
						{
							Date: 1738713236,
							Kwh:  4.5,
						},
						{
							Date: 1738716836,
							Kwh:  9.7,
						},
						{
							Date: 1738720436,
							Kwh:  1.2,
						},
						{
							Date: 1738724036,
							Kwh:  5.6,
						},
						{
							Date: 1738727636,
							Kwh:  7.8,
						},
						{
							Date: 1738731236,
							Kwh:  3.5,
						},
						{
							Date: 1738734836,
							Kwh:  2.4,
						},
						{
							Date: 1738738436,
							Kwh:  8.9,
						},
						{
							Date: 1738742036,
							Kwh:  6.9,
						},
						{
							Date: 1738745636,
							Kwh:  4.1,
						},
						{
							Date: 1738749236,
							Kwh:  9.9,
						},
						{
							Date: 1738752836,
							Kwh:  1.1,
						},
						{
							Date: 1738756436,
							Kwh:  2.3,
						},
						{
							Date: 1738760036,
							Kwh:  7.8,
						},
						{
							Date: 1738763636,
							Kwh:  4.5,
						},
						{
							Date: 1738767236,
							Kwh:  9.2,
						},
						{
							Date: 1738770836,
							Kwh:  1.6,
						},
						{
							Date: 1738774436,
							Kwh:  5.7,
						},
						{
							Date: 1738778036,
							Kwh:  8.4,
						},
						{
							Date: 1738781636,
							Kwh:  3.9,
						},
						{
							Date: 1738785236,
							Kwh:  6.1,
						},
						{
							Date: 1738788836,
							Kwh:  2.8,
						},
						{
							Date: 1738792436,
							Kwh:  7.3,
						},
						{
							Date: 1738796036,
							Kwh:  4.9,
						},
						{
							Date: 1738799636,
							Kwh:  9.6,
						},
						{
							Date: 1738803236,
							Kwh:  1.4,
						},
						{
							Date: 1738806836,
							Kwh:  5.2,
						},
						{
							Date: 1738810436,
							Kwh:  8.7,
						},
						{
							Date: 1738814036,
							Kwh:  3.3,
						},
						{
							Date: 1738817636,
							Kwh:  6.8,
						},
						{
							Date: 1738821236,
							Kwh:  2.5,
						},
						{
							Date: 1738824836,
							Kwh:  7.6,
						},
						{
							Date: 1738828436,
							Kwh:  4.2,
						},
						{
							Date: 1738832036,
							Kwh:  9.9,
						},
						{
							Date: 1738835636,
							Kwh:  1.8,
						},
						{
							Date: 1738839236,
							Kwh:  5.5,
						},
						{
							Date: 1738842836,
							Kwh:  8.1,
						},
						{
							Date: 1738846436,
							Kwh:  3.7,
						},
						{
							Date: 1738850036,
							Kwh:  6.4,
						},
						{
							Date: 1738853636,
							Kwh:  2.9,
						},
						{
							Date: 1738857236,
							Kwh:  7.9,
						},
						{
							Date: 1738860836,
							Kwh:  4.6,
						},
						{
							Date: 1738864436,
							Kwh:  9.3,
						},
						{
							Date: 1738868036,
							Kwh:  1.2,
						},
						{
							Date: 1738871636,
							Kwh:  5.8,
						},
						{
							Date: 1738875236,
							Kwh:  8.5,
						},
						{
							Date: 1738878836,
							Kwh:  3.1,
						},
						{
							Date: 1738882436,
							Kwh:  6.7,
						},
						{
							Date: 1738886036,
							Kwh:  2.4,
						},
						{
							Date: 1738889636,
							Kwh:  7.4,
						},
						{
							Date: 1738893236,
							Kwh:  4.8,
						},
						{
							Date: 1738896836,
							Kwh:  9.7,
						},
						{
							Date: 1738900436,
							Kwh:  1.5,
						},
						{
							Date: 1738904036,
							Kwh:  5.3,
						},
						{
							Date: 1738907636,
							Kwh:  8.9,
						},
						{
							Date: 1738911236,
							Kwh:  3.6,
						},
						{
							Date: 1738914836,
							Kwh:  6.2,
						},
						{
							Date: 1738918436,
							Kwh:  2.7,
						},
						{
							Date: 1738922036,
							Kwh:  7.1,
						},
						{
							Date: 1738925636,
							Kwh:  4.4,
						},
						{
							Date: 1738929236,
							Kwh:  5.9,
						},
						{
							Date: 1738932836,
							Kwh:  8.3,
						},
						{
							Date: 1738936436,
							Kwh:  3.4,
						},
						{
							Date: 1738940036,
							Kwh:  6.6,
						},
						{
							Date: 1738943636,
							Kwh:  2.2,
						},
						{
							Date: 1738947236,
							Kwh:  7.7,
						},
						{
							Date: 1738950836,
							Kwh:  4.3,
						},
						{
							Date: 1738954436,
							Kwh:  9.5,
						},
						{
							Date: 1738958036,
							Kwh:  1.3,
						},
						{
							Date: 1738961636,
							Kwh:  5.6,
						},
						{
							Date: 1738965236,
							Kwh:  8.8,
						},
						{
							Date: 1738968836,
							Kwh:  3.8,
						},
						{
							Date: 1738972436,
							Kwh:  6.9,
						},
						{
							Date: 1738976036,
							Kwh:  2.5,
						},
						{
							Date: 1738979636,
							Kwh:  7.2,
						},
						{
							Date: 1738983236,
							Kwh:  4.7,
						},
						{
							Date: 1738986836,
							Kwh:  9.1,
						},
						{
							Date: 1738990436,
							Kwh:  1.7,
						},
						{
							Date: 1738994036,
							Kwh:  5.4,
						},
						{
							Date: 1738997636,
							Kwh:  8.6,
						},
						{
							Date: 1739001236,
							Kwh:  3.2,
						},
						{
							Date: 1739004836,
							Kwh:  6.3,
						},
						{
							Date: 1739008436,
							Kwh:  2.9,
						},
						{
							Date: 1739012036,
							Kwh:  7.5,
						},
						{
							Date: 1739015636,
							Kwh:  4.1,
						},
						{
							Date: 1739019236,
							Kwh:  9.8,
						},
						{
							Date: 1739022836,
							Kwh:  1.9,
						},
						{
							Date: 1739026436,
							Kwh:  5.8,
						},
						{
							Date: 1739030036,
							Kwh:  8.2,
						},
						{
							Date: 1739033636,
							Kwh:  3.6,
						},
						{
							Date: 1739037236,
							Kwh:  6.7,
						},
						{
							Date: 1739040836,
							Kwh:  2.4,
						},
						{
							Date: 1739044436,
							Kwh:  7.9,
						},
						{
							Date: 1739048036,
							Kwh:  4.5,
						},
						{
							Date: 1739051636,
							Kwh:  9.3,
						},
						{
							Date: 1739055236,
							Kwh:  1.2,
						},
						{
							Date: 1739058836,
							Kwh:  5.5,
						},
						{
							Date: 1739062436,
							Kwh:  8.4,
						},
						{
							Date: 1739066036,
							Kwh:  3.7,
						},
						{
							Date: 1739069636,
							Kwh:  6.1,
						},
						{
							Date: 1739073236,
							Kwh:  2.8,
						},
						{
							Date: 1739076836,
							Kwh:  7.3,
						},
						{
							Date: 1739080436,
							Kwh:  4.9,
						},
						{
							Date: 1739084036,
							Kwh:  9.6,
						},
						{
							Date: 1739087636,
							Kwh:  1.4,
						},
						{
							Date: 1739091236,
							Kwh:  5.2,
						},
						{
							Date: 1739094836,
							Kwh:  8.7,
						},
						{
							Date: 1739098436,
							Kwh:  3.3,
						},
						{
							Date: 1739102036,
							Kwh:  6.8,
						},
						{
							Date: 1739105636,
							Kwh:  2.5,
						},
						{
							Date: 1739109236,
							Kwh:  7.6,
						},
						{
							Date: 1739112836,
							Kwh:  4.2,
						},
						{
							Date: 1739116436,
							Kwh:  9.9,
						},
						{
							Date: 1739120036,
							Kwh:  1.8,
						},
						{
							Date: 1739123636,
							Kwh:  5.5,
						},
						{
							Date: 1739127236,
							Kwh:  8.1,
						},
						{
							Date: 1739130836,
							Kwh:  3.7,
						},
						{
							Date: 1739134436,
							Kwh:  6.4,
						},
						{
							Date: 1739138036,
							Kwh:  2.9,
						},
						{
							Date: 1739141636,
							Kwh:  7.9,
						},
						{
							Date: 1739145236,
							Kwh:  4.6,
						},
						{
							Date: 1739148836,
							Kwh:  9.3,
						},
						{
							Date: 1739152436,
							Kwh:  1.2,
						},
						{
							Date: 1739156036,
							Kwh:  5.8,
						},
						{
							Date: 1739159636,
							Kwh:  8.5,
						},
						{
							Date: 1739163236,
							Kwh:  3.1,
						},
						{
							Date: 1739166836,
							Kwh:  6.7,
						},
						{
							Date: 1739170436,
							Kwh:  2.4,
						},
						{
							Date: 1739174036,
							Kwh:  7.4,
						},
						{
							Date: 1739177636,
							Kwh:  4.8,
						},
						{
							Date: 1739181236,
							Kwh:  9.7,
						},
						{
							Date: 1739184836,
							Kwh:  1.5,
						},
						{
							Date: 1739188436,
							Kwh:  5.3,
						},
						{
							Date: 1739192036,
							Kwh:  8.9,
						},
						{
							Date: 1739195636,
							Kwh:  3.6,
						},
						{
							Date: 1739199236,
							Kwh:  6.2,
						},
						{
							Date: 1739202836,
							Kwh:  2.7,
						},
						{
							Date: 1739206436,
							Kwh:  7.1,
						},
						{
							Date: 1739210036,
							Kwh:  4.4,
						},
						{
							Date: 1739213636,
							Kwh:  5.9,
						},
						{
							Date: 1739217236,
							Kwh:  8.3,
						},
						{
							Date: 1739220836,
							Kwh:  3.4,
						},
						{
							Date: 1739224436,
							Kwh:  6.6,
						},
						{
							Date: 1739228036,
							Kwh:  2.2,
						},
						{
							Date: 1739231636,
							Kwh:  7.7,
						},
						{
							Date: 1739235236,
							Kwh:  4.3,
						},
						{
							Date: 1739238836,
							Kwh:  9.5,
						},
						{
							Date: 1739242436,
							Kwh:  1.3,
						},
						{
							Date: 1739246036,
							Kwh:  5.6,
						},
						{
							Date: 1739249636,
							Kwh:  8.8,
						},
						{
							Date: 1739253236,
							Kwh:  3.8,
						},
						{
							Date: 1739256836,
							Kwh:  6.9,
						},
						{
							Date: 1739260436,
							Kwh:  2.5,
						},
						{
							Date: 1739264036,
							Kwh:  7.2,
						},
						{
							Date: 1739267636,
							Kwh:  4.7,
						},
						{
							Date: 1739271236,
							Kwh:  9.1,
						},
						{
							Date: 1739274836,
							Kwh:  1.7,
						},
						{
							Date: 1739278436,
							Kwh:  5.4,
						},
						{
							Date: 1739282036,
							Kwh:  8.6,
						},
						{
							Date: 1739285636,
							Kwh:  3.2,
						},
						{
							Date: 1739289236,
							Kwh:  6.3,
						},
						{
							Date: 1739292836,
							Kwh:  2.9,
						},
						{
							Date: 1739296436,
							Kwh:  7.5,
						},
						{
							Date: 1739300036,
							Kwh:  4.1,
						},
						{
							Date: 1739303636,
							Kwh:  9.8,
						},
						{
							Date: 1739307236,
							Kwh:  1.9,
						},
						{
							Date: 1739310836,
							Kwh:  5.8,
						},
						{
							Date: 1739314436,
							Kwh:  8.2,
						},
						{
							Date: 1739318036,
							Kwh:  3.6,
						},
						{
							Date: 1739321636,
							Kwh:  6.7,
						},
						{
							Date: 1739325236,
							Kwh:  2.4,
						},
						{
							Date: 1739328836,
							Kwh:  7.9,
						},
						{
							Date: 1739332436,
							Kwh:  4.5,
						},
						{
							Date: 1739336036,
							Kwh:  9.3,
						},
						{
							Date: 1739339636,
							Kwh:  1.2,
						},
						{
							Date: 1739343236,
							Kwh:  5.5,
						},
						{
							Date: 1739346836,
							Kwh:  8.4,
						},
						{
							Date: 1739350436,
							Kwh:  3.7,
						},
						{
							Date: 1739354036,
							Kwh:  6.1,
						},
						{
							Date: 1739357636,
							Kwh:  2.8,
						},
						{
							Date: 1739361236,
							Kwh:  7.3,
						},
						{
							Date: 1739364836,
							Kwh:  4.9,
						},
						{
							Date: 1739368436,
							Kwh:  9.6,
						},
						{
							Date: 1739372036,
							Kwh:  1.4,
						},
						{
							Date: 1739375636,
							Kwh:  5.2,
						},
						{
							Date: 1739379236,
							Kwh:  8.7,
						},
						{
							Date: 1739382836,
							Kwh:  3.3,
						},
						{
							Date: 1739386436,
							Kwh:  6.8,
						},
						{
							Date: 1739390036,
							Kwh:  2.5,
						},
						{
							Date: 1739393636,
							Kwh:  7.6,
						},
						{
							Date: 1739397236,
							Kwh:  4.2,
						},
						{
							Date: 1739400836,
							Kwh:  9.9,
						},
						{
							Date: 1739404436,
							Kwh:  1.8,
						},
						{
							Date: 1739408036,
							Kwh:  5.5,
						},
						{
							Date: 1739411636,
							Kwh:  8.1,
						},
						{
							Date: 1739415236,
							Kwh:  3.7,
						},
						{
							Date: 1739418836,
							Kwh:  6.4,
						},
						{
							Date: 1739422436,
							Kwh:  2.9,
						},
						{
							Date: 1739426036,
							Kwh:  7.9,
						},
						{
							Date: 1739429636,
							Kwh:  4.6,
						},
						{
							Date: 1739433236,
							Kwh:  9.3,
						},
						{
							Date: 1739436836,
							Kwh:  1.2,
						},
						{
							Date: 1739440436,
							Kwh:  5.8,
						},
						{
							Date: 1739444036,
							Kwh:  8.5,
						},
						{
							Date: 1739447636,
							Kwh:  3.1,
						},
						{
							Date: 1739451236,
							Kwh:  6.7,
						},
						{
							Date: 1739454836,
							Kwh:  2.4,
						},
						{
							Date: 1739458436,
							Kwh:  7.4,
						},
						{
							Date: 1739462036,
							Kwh:  4.8,
						},
						{
							Date: 1739465636,
							Kwh:  9.7,
						},
						{
							Date: 1739469236,
							Kwh:  1.5,
						},
						{
							Date: 1739472836,
							Kwh:  5.3,
						},
						{
							Date: 1739476436,
							Kwh:  8.9,
						},
						{
							Date: 1739480036,
							Kwh:  3.6,
						},
						{
							Date: 1739483636,
							Kwh:  6.2,
						},
						{
							Date: 1739487236,
							Kwh:  2.7,
						},
						{
							Date: 1739490836,
							Kwh:  7.1,
						},
						{
							Date: 1739494436,
							Kwh:  4.4,
						},
						{
							Date: 1739498036,
							Kwh:  5.9,
						},
						{
							Date: 1739501636,
							Kwh:  8.3,
						},
						{
							Date: 1739505236,
							Kwh:  3.4,
						},
						{
							Date: 1739508836,
							Kwh:  6.6,
						},
						{
							Date: 1739512436,
							Kwh:  2.2,
						},
						{
							Date: 1739516036,
							Kwh:  7.7,
						},
						{
							Date: 1739519636,
							Kwh:  4.3,
						},
						{
							Date: 1739523236,
							Kwh:  9.5,
						},
						{
							Date: 1739526836,
							Kwh:  1.3,
						},
						{
							Date: 1739530436,
							Kwh:  5.6,
						},
						{
							Date: 1739534036,
							Kwh:  8.8,
						},
						{
							Date: 1739537636,
							Kwh:  3.8,
						},
						{
							Date: 1739541236,
							Kwh:  6.9,
						},
						{
							Date: 1739544836,
							Kwh:  2.5,
						},
						{
							Date: 1739548436,
							Kwh:  7.2,
						},
						{
							Date: 1739552036,
							Kwh:  4.7,
						},
						{
							Date: 1739555636,
							Kwh:  9.1,
						},
						{
							Date: 1739559236,
							Kwh:  1.7,
						},
						{
							Date: 1739562836,
							Kwh:  5.4,
						},
						{
							Date: 1739566436,
							Kwh:  8.6,
						},
						{
							Date: 1739570036,
							Kwh:  3.2,
						},
						{
							Date: 1739573636,
							Kwh:  6.3,
						},
						{
							Date: 1739577236,
							Kwh:  2.9,
						},
						{
							Date: 1739580836,
							Kwh:  7.5,
						},
						{
							Date: 1739584436,
							Kwh:  4.1,
						},
						{
							Date: 1739588036,
							Kwh:  9.8,
						},
						{
							Date: 1739591636,
							Kwh:  1.9,
						},
						{
							Date: 1739595236,
							Kwh:  5.8,
						},
						{
							Date: 1739598836,
							Kwh:  8.2,
						},
						{
							Date: 1739602436,
							Kwh:  3.6,
						},
						{
							Date: 1739606036,
							Kwh:  6.7,
						},
						{
							Date: 1739609636,
							Kwh:  2.4,
						},
						{
							Date: 1739613236,
							Kwh:  7.9,
						},
						{
							Date: 1739616836,
							Kwh:  4.5,
						},
						{
							Date: 1739620436,
							Kwh:  9.3,
						},
						{
							Date: 1739624036,
							Kwh:  1.2,
						},
						{
							Date: 1739627636,
							Kwh:  5.5,
						},
						{
							Date: 1739631236,
							Kwh:  8.4,
						},
						{
							Date: 1739634836,
							Kwh:  3.7,
						},
						{
							Date: 1739638436,
							Kwh:  6.1,
						},
						{
							Date: 1739642036,
							Kwh:  2.8,
						},
						{
							Date: 1739645636,
							Kwh:  7.3,
						},
						{
							Date: 1739649236,
							Kwh:  4.9,
						},
						{
							Date: 1739652836,
							Kwh:  9.6,
						},
						{
							Date: 1739656436,
							Kwh:  1.4,
						},
						{
							Date: 1739660036,
							Kwh:  5.2,
						},
						{
							Date: 1739663636,
							Kwh:  8.7,
						},
						{
							Date: 1739667236,
							Kwh:  3.3,
						},
						{
							Date: 1739670836,
							Kwh:  6.8,
						},
						{
							Date: 1739674436,
							Kwh:  2.5,
						},
						{
							Date: 1739678036,
							Kwh:  7.6,
						},
						{
							Date: 1739681636,
							Kwh:  4.2,
						},
						{
							Date: 1739685236,
							Kwh:  9.9,
						},
						{
							Date: 1739688836,
							Kwh:  1.8,
						},
						{
							Date: 1739692436,
							Kwh:  5.5,
						},
						{
							Date: 1739696036,
							Kwh:  8.1,
						},
						{
							Date: 1739699636,
							Kwh:  3.7,
						},
						{
							Date: 1739703236,
							Kwh:  6.4,
						},
						{
							Date: 1739706836,
							Kwh:  2.9,
						},
						{
							Date: 1739710436,
							Kwh:  7.9,
						},
						{
							Date: 1739714036,
							Kwh:  4.6,
						},
						{
							Date: 1739717636,
							Kwh:  9.3,
						},
						{
							Date: 1739721236,
							Kwh:  1.2,
						},
						{
							Date: 1739724836,
							Kwh:  5.8,
						},
						{
							Date: 1739728436,
							Kwh:  8.5,
						},
						{
							Date: 1739732036,
							Kwh:  3.1,
						},
						{
							Date: 1739735636,
							Kwh:  6.7,
						},
						{
							Date: 1739739236,
							Kwh:  2.4,
						},
						{
							Date: 1739742836,
							Kwh:  7.4,
						},
						{
							Date: 1739746436,
							Kwh:  4.8,
						},
						{
							Date: 1739750036,
							Kwh:  9.7,
						},
						{
							Date: 1739753636,
							Kwh:  1.5,
						},
						{
							Date: 1739757236,
							Kwh:  5.3,
						},
						{
							Date: 1739760836,
							Kwh:  8.9,
						},
						{
							Date: 1739764436,
							Kwh:  3.6,
						},
						{
							Date: 1739768036,
							Kwh:  6.2,
						},
						{
							Date: 1739771636,
							Kwh:  2.7,
						},
						{
							Date: 1739775236,
							Kwh:  7.1,
						},
						{
							Date: 1739778836,
							Kwh:  4.4,
						},
						{
							Date: 1739782436,
							Kwh:  5.9,
						},
						{
							Date: 1739786036,
							Kwh:  8.3,
						},
						{
							Date: 1739789636,
							Kwh:  3.4,
						},
						{
							Date: 1739793236,
							Kwh:  6.6,
						},
						{
							Date: 1739796836,
							Kwh:  2.2,
						},
						{
							Date: 1739800436,
							Kwh:  7.7,
						},
						{
							Date: 1739804036,
							Kwh:  4.3,
						},
						{
							Date: 1739807636,
							Kwh:  9.5,
						},
						{
							Date: 1739811236,
							Kwh:  1.3,
						},
						{
							Date: 1739814836,
							Kwh:  5.6,
						},
						{
							Date: 1739818436,
							Kwh:  8.8,
						},
						{
							Date: 1739822036,
							Kwh:  3.8,
						},
						{
							Date: 1739825636,
							Kwh:  6.9,
						},
						{
							Date: 1739829236,
							Kwh:  2.5,
						},
						{
							Date: 1739832836,
							Kwh:  7.2,
						},
						{
							Date: 1739836436,
							Kwh:  4.7,
						},
						{
							Date: 1739840036,
							Kwh:  9.1,
						},
						{
							Date: 1739843636,
							Kwh:  1.7,
						},
						{
							Date: 1739847236,
							Kwh:  5.4,
						},
						{
							Date: 1739850836,
							Kwh:  8.6,
						},
						{
							Date: 1739854436,
							Kwh:  3.2,
						},
						{
							Date: 1739858036,
							Kwh:  6.3,
						},
						{
							Date: 1739861636,
							Kwh:  2.9,
						},
						{
							Date: 1739865236,
							Kwh:  7.5,
						},
						{
							Date: 1739868836,
							Kwh:  4.1,
						},
						{
							Date: 1739872436,
							Kwh:  9.8,
						},
						{
							Date: 1739876036,
							Kwh:  1.9,
						},
						{
							Date: 1739879636,
							Kwh:  5.8,
						},
						{
							Date: 1739883236,
							Kwh:  8.2,
						},
						{
							Date: 1739886836,
							Kwh:  3.6,
						},
						{
							Date: 1739890436,
							Kwh:  6.7,
						},
						{
							Date: 1739894036,
							Kwh:  2.4,
						},
						{
							Date: 1739897636,
							Kwh:  7.9,
						},
						{
							Date: 1739901236,
							Kwh:  4.5,
						},
						{
							Date: 1739904836,
							Kwh:  9.3,
						},
						{
							Date: 1739908436,
							Kwh:  1.2,
						},
						{
							Date: 1739912036,
							Kwh:  5.5,
						},
						{
							Date: 1739915636,
							Kwh:  8.4,
						},
						{
							Date: 1739919236,
							Kwh:  3.7,
						},
						{
							Date: 1739922836,
							Kwh:  6.1,
						},
						{
							Date: 1739926436,
							Kwh:  2.8,
						},
						{
							Date: 1739930036,
							Kwh:  7.3,
						},
						{
							Date: 1739933636,
							Kwh:  4.9,
						},
						{
							Date: 1739937236,
							Kwh:  9.6,
						},
						{
							Date: 1739940836,
							Kwh:  1.4,
						},
						{
							Date: 1739944436,
							Kwh:  5.2,
						},
						{
							Date: 1739948036,
							Kwh:  8.7,
						},
						{
							Date: 1739951636,
							Kwh:  3.3,
						},
						{
							Date: 1739955236,
							Kwh:  6.8,
						},
						{
							Date: 1739958836,
							Kwh:  2.5,
						},
						{
							Date: 1739962436,
							Kwh:  7.6,
						},
						{
							Date: 1739966036,
							Kwh:  4.2,
						},
						{
							Date: 1739969636,
							Kwh:  9.9,
						},
						{
							Date: 1739973236,
							Kwh:  1.8,
						},
						{
							Date: 1739976836,
							Kwh:  5.5,
						},
						{
							Date: 1739980436,
							Kwh:  8.1,
						},
						{
							Date: 1739984036,
							Kwh:  3.7,
						},
						{
							Date: 1739987636,
							Kwh:  6.4,
						},
						{
							Date: 1739991236,
							Kwh:  2.9,
						},
						{
							Date: 1739994836,
							Kwh:  7.9,
						},
						{
							Date: 1739998436,
							Kwh:  4.6,
						},
						{
							Date: 1740002036,
							Kwh:  9.3,
						},
						{
							Date: 1740005636,
							Kwh:  1.2,
						},
						{
							Date: 1740009236,
							Kwh:  5.8,
						},
						{
							Date: 1740012836,
							Kwh:  8.5,
						},
						{
							Date: 1740016436,
							Kwh:  3.1,
						},
						{
							Date: 1740020036,
							Kwh:  6.7,
						},
						{
							Date: 1740023636,
							Kwh:  2.4,
						},
						{
							Date: 1740027236,
							Kwh:  7.4,
						},
						{
							Date: 1740030836,
							Kwh:  4.8,
						},
						{
							Date: 1740034436,
							Kwh:  9.7,
						},
						{
							Date: 1740038036,
							Kwh:  1.5,
						},
						{
							Date: 1740041636,
							Kwh:  5.3,
						},
						{
							Date: 1740045236,
							Kwh:  8.9,
						},
						{
							Date: 1740048836,
							Kwh:  3.6,
						},
						{
							Date: 1740052436,
							Kwh:  6.2,
						},
						{
							Date: 1740056036,
							Kwh:  2.7,
						},
						{
							Date: 1740059636,
							Kwh:  7.1,
						},
						{
							Date: 1740063236,
							Kwh:  4.4,
						},
						{
							Date: 1740066836,
							Kwh:  5.9,
						},
						{
							Date: 1740070436,
							Kwh:  8.3,
						},
						{
							Date: 1740074036,
							Kwh:  3.4,
						},
						{
							Date: 1740077636,
							Kwh:  6.6,
						},
						{
							Date: 1740081236,
							Kwh:  2.2,
						},
						{
							Date: 1740084836,
							Kwh:  7.7,
						},
						{
							Date: 1740088436,
							Kwh:  4.3,
						},
						{
							Date: 1740092036,
							Kwh:  9.5,
						},
						{
							Date: 1740095636,
							Kwh:  1.3,
						},
						{
							Date: 1740099236,
							Kwh:  5.6,
						},
						{
							Date: 1740102836,
							Kwh:  8.8,
						},
						{
							Date: 1740106436,
							Kwh:  3.8,
						},
						{
							Date: 1740110036,
							Kwh:  6.9,
						},
						{
							Date: 1740113636,
							Kwh:  2.5,
						},
						{
							Date: 1740117236,
							Kwh:  7.2,
						},
						{
							Date: 1740120836,
							Kwh:  4.7,
						},
						{
							Date: 1740124436,
							Kwh:  9.1,
						},
						{
							Date: 1740128036,
							Kwh:  1.7,
						},
						{
							Date: 1740131636,
							Kwh:  5.4,
						},
						{
							Date: 1740135236,
							Kwh:  8.6,
						},
						{
							Date: 1740138836,
							Kwh:  3.2,
						},
						{
							Date: 1740142436,
							Kwh:  6.3,
						},
						{
							Date: 1740146036,
							Kwh:  2.9,
						},
						{
							Date: 1740149636,
							Kwh:  7.5,
						},
						{
							Date: 1740153236,
							Kwh:  4.1,
						},
						{
							Date: 1740156836,
							Kwh:  9.8,
						},
						{
							Date: 1740160436,
							Kwh:  1.9,
						},
						{
							Date: 1740164036,
							Kwh:  5.8,
						},
						{
							Date: 1740167636,
							Kwh:  8.2,
						},
						{
							Date: 1740171236,
							Kwh:  3.6,
						},
						{
							Date: 1740174836,
							Kwh:  6.7,
						},
						{
							Date: 1740178436,
							Kwh:  2.4,
						},
						{
							Date: 1740182036,
							Kwh:  7.9,
						},
						{
							Date: 1740185636,
							Kwh:  4.5,
						},
						{
							Date: 1740189236,
							Kwh:  9.3,
						},
						{
							Date: 1740192836,
							Kwh:  1.2,
						},
						{
							Date: 1740196436,
							Kwh:  5.5,
						},
						{
							Date: 1740200036,
							Kwh:  8.4,
						},
						{
							Date: 1740203636,
							Kwh:  3.7,
						},
						{
							Date: 1740207236,
							Kwh:  6.1,
						},
						{
							Date: 1740210836,
							Kwh:  2.8,
						},
						{
							Date: 1740214436,
							Kwh:  7.3,
						},
						{
							Date: 1740218036,
							Kwh:  4.9,
						},
						{
							Date: 1740221636,
							Kwh:  9.6,
						},
						{
							Date: 1740225236,
							Kwh:  1.4,
						},
						{
							Date: 1740228836,
							Kwh:  5.2,
						},
						{
							Date: 1740232436,
							Kwh:  8.7,
						},
						{
							Date: 1740236036,
							Kwh:  3.3,
						},
						{
							Date: 1740239636,
							Kwh:  6.8,
						},
						{
							Date: 1740243236,
							Kwh:  2.5,
						},
						{
							Date: 1740246836,
							Kwh:  7.6,
						},
						{
							Date: 1740250436,
							Kwh:  4.2,
						},
						{
							Date: 1740254036,
							Kwh:  9.9,
						},
						{
							Date: 1740257636,
							Kwh:  1.8,
						},
						{
							Date: 1740261236,
							Kwh:  5.5,
						},
						{
							Date: 1740264836,
							Kwh:  8.1,
						},
						{
							Date: 1740268436,
							Kwh:  3.7,
						},
						{
							Date: 1740272036,
							Kwh:  6.4,
						},
						{
							Date: 1740275636,
							Kwh:  2.9,
						},
						{
							Date: 1740279236,
							Kwh:  7.9,
						},
						{
							Date: 1740282836,
							Kwh:  4.6,
						},
						{
							Date: 1740286436,
							Kwh:  9.3,
						},
						{
							Date: 1740290036,
							Kwh:  1.2,
						},
						{
							Date: 1740293636,
							Kwh:  5.8,
						},
						{
							Date: 1740297236,
							Kwh:  8.5,
						},
						{
							Date: 1740300836,
							Kwh:  3.1,
						},
						{
							Date: 1740304436,
							Kwh:  6.7,
						},
						{
							Date: 1740308036,
							Kwh:  2.4,
						},
						{
							Date: 1740311636,
							Kwh:  7.4,
						},
						{
							Date: 1740315236,
							Kwh:  4.8,
						},
						{
							Date: 1740318836,
							Kwh:  9.7,
						},
						{
							Date: 1740322436,
							Kwh:  1.5,
						},
						{
							Date: 1740326036,
							Kwh:  5.3,
						},
						{
							Date: 1740329636,
							Kwh:  8.9,
						},
						{
							Date: 1740333236,
							Kwh:  3.6,
						},
						{
							Date: 1740336836,
							Kwh:  6.2,
						},
						{
							Date: 1740340436,
							Kwh:  2.7,
						},
						{
							Date: 1740344036,
							Kwh:  7.1,
						},
						{
							Date: 1740347636,
							Kwh:  4.4,
						},
						{
							Date: 1740351236,
							Kwh:  5.9,
						},
						{
							Date: 1740354836,
							Kwh:  8.3,
						},
						{
							Date: 1740358436,
							Kwh:  3.4,
						},
						{
							Date: 1740362036,
							Kwh:  6.6,
						},
						{
							Date: 1740365636,
							Kwh:  2.2,
						},
						{
							Date: 1740369236,
							Kwh:  7.7,
						},
						{
							Date: 1740372836,
							Kwh:  4.3,
						},
						{
							Date: 1740376436,
							Kwh:  9.5,
						},
						{
							Date: 1740380036,
							Kwh:  1.3,
						},
						{
							Date: 1740383636,
							Kwh:  5.6,
						},
						{
							Date: 1740387236,
							Kwh:  8.8,
						},
						{
							Date: 1740390836,
							Kwh:  3.8,
						},
						{
							Date: 1740394436,
							Kwh:  6.9,
						},
						{
							Date: 1740398036,
							Kwh:  2.5,
						},
						{
							Date: 1740401636,
							Kwh:  7.2,
						},
						{
							Date: 1740405236,
							Kwh:  4.7,
						},
						{
							Date: 1740408836,
							Kwh:  9.1,
						},
						{
							Date: 1740412436,
							Kwh:  1.7,
						},
						{
							Date: 1740416036,
							Kwh:  5.4,
						},
						{
							Date: 1740419636,
							Kwh:  8.6,
						},
						{
							Date: 1740423236,
							Kwh:  3.2,
						},
						{
							Date: 1740426836,
							Kwh:  6.3,
						},
						{
							Date: 1740430436,
							Kwh:  2.9,
						},
						{
							Date: 1740434036,
							Kwh:  7.5,
						},
						{
							Date: 1740437636,
							Kwh:  4.1,
						},
						{
							Date: 1740441236,
							Kwh:  9.8,
						},
						{
							Date: 1740444836,
							Kwh:  1.9,
						},
						{
							Date: 1740448436,
							Kwh:  5.8,
						},
						{
							Date: 1740452036,
							Kwh:  8.2,
						},
						{
							Date: 1740455636,
							Kwh:  3.6,
						},
						{
							Date: 1740459236,
							Kwh:  6.7,
						},
						{
							Date: 1740462836,
							Kwh:  2.4,
						},
						{
							Date: 1740466436,
							Kwh:  7.9,
						},
						{
							Date: 1740470036,
							Kwh:  4.5,
						},
						{
							Date: 1740473636,
							Kwh:  9.3,
						},
						{
							Date: 1740477236,
							Kwh:  1.2,
						},
						{
							Date: 1740480836,
							Kwh:  5.5,
						},
						{
							Date: 1740484436,
							Kwh:  8.4,
						},
						{
							Date: 1740488036,
							Kwh:  3.7,
						},
						{
							Date: 1740491636,
							Kwh:  6.1,
						},
						{
							Date: 1740495236,
							Kwh:  2.8,
						},
						{
							Date: 1740498836,
							Kwh:  7.3,
						},
						{
							Date: 1740502436,
							Kwh:  4.9,
						},
						{
							Date: 1740506036,
							Kwh:  9.6,
						},
						{
							Date: 1740509636,
							Kwh:  1.4,
						},
						{
							Date: 1740513236,
							Kwh:  5.2,
						},
						{
							Date: 1740516836,
							Kwh:  8.7,
						},
						{
							Date: 1740520436,
							Kwh:  3.3,
						},
						{
							Date: 1740524036,
							Kwh:  6.8,
						},
						{
							Date: 1740527636,
							Kwh:  2.5,
						},
						{
							Date: 1740531236,
							Kwh:  7.6,
						},
						{
							Date: 1740534836,
							Kwh:  4.2,
						},
						{
							Date: 1740538436,
							Kwh:  9.9,
						},
						{
							Date: 1740542036,
							Kwh:  1.8,
						},
						{
							Date: 1740545636,
							Kwh:  5.5,
						},
						{
							Date: 1740549236,
							Kwh:  8.1,
						},
						{
							Date: 1740552836,
							Kwh:  3.7,
						},
						{
							Date: 1740556436,
							Kwh:  6.4,
						},
						{
							Date: 1740560036,
							Kwh:  2.9,
						},
						{
							Date: 1740563636,
							Kwh:  7.9,
						},
						{
							Date: 1740567236,
							Kwh:  4.6,
						},
						{
							Date: 1740570836,
							Kwh:  9.3,
						},
						{
							Date: 1740574436,
							Kwh:  1.2,
						},
						{
							Date: 1740578036,
							Kwh:  5.8,
						},
						{
							Date: 1740581636,
							Kwh:  8.5,
						},
						{
							Date: 1740585236,
							Kwh:  3.1,
						},
						{
							Date: 1740588836,
							Kwh:  6.7,
						},
						{
							Date: 1740592436,
							Kwh:  2.4,
						},
						{
							Date: 1740596036,
							Kwh:  7.4,
						},
						{
							Date: 1740599636,
							Kwh:  4.8,
						},
						{
							Date: 1740603236,
							Kwh:  9.7,
						},
						{
							Date: 1740606836,
							Kwh:  1.5,
						},
						{
							Date: 1740610436,
							Kwh:  5.3,
						},
						{
							Date: 1740614036,
							Kwh:  8.9,
						},
						{
							Date: 1740617636,
							Kwh:  3.6,
						},
						{
							Date: 1740621236,
							Kwh:  6.2,
						},
						{
							Date: 1740624836,
							Kwh:  2.7,
						},
						{
							Date: 1740628436,
							Kwh:  7.1,
						},
						{
							Date: 1740632036,
							Kwh:  4.4,
						},
						{
							Date: 1740635636,
							Kwh:  5.9,
						},
						{
							Date: 1740639236,
							Kwh:  8.3,
						},
						{
							Date: 1740642836,
							Kwh:  3.4,
						},
						{
							Date: 1740646436,
							Kwh:  6.6,
						},
						{
							Date: 1740650036,
							Kwh:  2.2,
						},
						{
							Date: 1740653636,
							Kwh:  7.7,
						},
						{
							Date: 1740657236,
							Kwh:  4.3,
						},
						{
							Date: 1740660836,
							Kwh:  9.5,
						},
						{
							Date: 1740664436,
							Kwh:  1.3,
						},
						{
							Date: 1740668036,
							Kwh:  5.6,
						},
						{
							Date: 1740671636,
							Kwh:  8.8,
						},
						{
							Date: 1740675236,
							Kwh:  3.8,
						},
						{
							Date: 1740678836,
							Kwh:  6.9,
						},
						{
							Date: 1740682436,
							Kwh:  2.5,
						},
						{
							Date: 1740686036,
							Kwh:  7.2,
						},
						{
							Date: 1740689636,
							Kwh:  4.7,
						},
						{
							Date: 1740693236,
							Kwh:  9.1,
						},
						{
							Date: 1740696836,
							Kwh:  1.7,
						},
						{
							Date: 1740700436,
							Kwh:  5.4,
						},
						{
							Date: 1740704036,
							Kwh:  8.6,
						},
						{
							Date: 1740707636,
							Kwh:  3.2,
						},
						{
							Date: 1740711236,
							Kwh:  6.3,
						},
						{
							Date: 1740714836,
							Kwh:  2.9,
						},
						{
							Date: 1740718436,
							Kwh:  7.5,
						},
						{
							Date: 1740722036,
							Kwh:  4.1,
						},
						{
							Date: 1740725636,
							Kwh:  9.8,
						},
						{
							Date: 1740729236,
							Kwh:  1.9,
						},
						{
							Date: 1740732836,
							Kwh:  5.8,
						},
						{
							Date: 1740736436,
							Kwh:  8.2,
						},
						{
							Date: 1740740036,
							Kwh:  3.6,
						},
						{
							Date: 1740743636,
							Kwh:  6.7,
						},
						{
							Date: 1740747236,
							Kwh:  2.4,
						},
						{
							Date: 1740750836,
							Kwh:  7.9,
						},
						{
							Date: 1740754436,
							Kwh:  4.5,
						},
						{
							Date: 1740758036,
							Kwh:  9.3,
						},
						{
							Date: 1740761636,
							Kwh:  1.2,
						},
						{
							Date: 1740765236,
							Kwh:  5.5,
						},
						{
							Date: 1740768836,
							Kwh:  8.4,
						},
						{
							Date: 1740772436,
							Kwh:  3.7,
						},
						{
							Date: 1740776036,
							Kwh:  6.1,
						},
						{
							Date: 1740779636,
							Kwh:  2.8,
						},
						{
							Date: 1740783236,
							Kwh:  7.3,
						},
						{
							Date: 1740786836,
							Kwh:  4.9,
						},
						{
							Date: 1740790436,
							Kwh:  9.6,
						},
						{
							Date: 1740794036,
							Kwh:  1.4,
						},
						{
							Date: 1740797636,
							Kwh:  5.2,
						},
						{
							Date: 1740801236,
							Kwh:  8.7,
						},
						{
							Date: 1740804836,
							Kwh:  3.3,
						},
						{
							Date: 1740808436,
							Kwh:  6.8,
						},
						{
							Date: 1740812036,
							Kwh:  2.5,
						},
						{
							Date: 1740815636,
							Kwh:  7.6,
						},
						{
							Date: 1740819236,
							Kwh:  4.2,
						},
						{
							Date: 1740822836,
							Kwh:  9.9,
						},
						{
							Date: 1740826436,
							Kwh:  1.8,
						},
						{
							Date: 1740830036,
							Kwh:  5.5,
						},
						{
							Date: 1740833636,
							Kwh:  8.1,
						},
						{
							Date: 1740837236,
							Kwh:  3.7,
						},
						{
							Date: 1740840836,
							Kwh:  6.4,
						},
						{
							Date: 1740844436,
							Kwh:  2.9,
						},
						{
							Date: 1740848036,
							Kwh:  7.9,
						},
						{
							Date: 1740851636,
							Kwh:  4.6,
						},
						{
							Date: 1740855236,
							Kwh:  9.3,
						},
						{
							Date: 1740858836,
							Kwh:  1.2,
						},
						{
							Date: 1740862436,
							Kwh:  5.8,
						},
						{
							Date: 1740866036,
							Kwh:  8.5,
						},
						{
							Date: 1740869636,
							Kwh:  3.1,
						},
						{
							Date: 1740873236,
							Kwh:  6.7,
						},
						{
							Date: 1740876836,
							Kwh:  2.4,
						},
						{
							Date: 1740880436,
							Kwh:  7.4,
						},
						{
							Date: 1740884036,
							Kwh:  4.8,
						},
						{
							Date: 1740887636,
							Kwh:  9.7,
						},
						{
							Date: 1740891236,
							Kwh:  1.5,
						},
						{
							Date: 1740894836,
							Kwh:  5.3,
						},
						{
							Date: 1740898436,
							Kwh:  8.9,
						},
						{
							Date: 1740902036,
							Kwh:  3.6,
						},
						{
							Date: 1740905636,
							Kwh:  6.2,
						},
						{
							Date: 1740909236,
							Kwh:  2.7,
						},
						{
							Date: 1740912836,
							Kwh:  7.1,
						},
						{
							Date: 1740916436,
							Kwh:  4.4,
						},
						{
							Date: 1740920036,
							Kwh:  5.9,
						},
						{
							Date: 1740923636,
							Kwh:  8.3,
						},
						{
							Date: 1740927236,
							Kwh:  3.4,
						},
						{
							Date: 1740930836,
							Kwh:  6.6,
						},
						{
							Date: 1740934436,
							Kwh:  2.2,
						},
						{
							Date: 1740938036,
							Kwh:  7.7,
						},
						{
							Date: 1740941636,
							Kwh:  4.3,
						},
						{
							Date: 1740945236,
							Kwh:  9.5,
						},
						{
							Date: 1740948836,
							Kwh:  1.3,
						},
						{
							Date: 1740952436,
							Kwh:  5.6,
						},
						{
							Date: 1740956036,
							Kwh:  8.8,
						},
						{
							Date: 1740959636,
							Kwh:  3.8,
						},
						{
							Date: 1740963236,
							Kwh:  6.9,
						},
						{
							Date: 1740966836,
							Kwh:  2.5,
						},
						{
							Date: 1740970436,
							Kwh:  7.2,
						},
						{
							Date: 1740974036,
							Kwh:  4.7,
						},
						{
							Date: 1740977636,
							Kwh:  9.1,
						},
						{
							Date: 1740981236,
							Kwh:  1.7,
						},
						{
							Date: 1740984836,
							Kwh:  5.4,
						},
						{
							Date: 1740988436,
							Kwh:  8.6,
						},
						{
							Date: 1740992036,
							Kwh:  3.2,
						},
						{
							Date: 1740995636,
							Kwh:  6.3,
						},
						{
							Date: 1740999236,
							Kwh:  2.9,
						},
						{
							Date: 1741002836,
							Kwh:  7.5,
						},
						{
							Date: 1741006436,
							Kwh:  4.1,
						},
						{
							Date: 1741010036,
							Kwh:  9.8,
						},
						{
							Date: 1741013636,
							Kwh:  1.9,
						},
						{
							Date: 1741017236,
							Kwh:  5.8,
						},
						{
							Date: 1741020836,
							Kwh:  8.2,
						},
						{
							Date: 1741024436,
							Kwh:  3.6,
						},
						{
							Date: 1741028036,
							Kwh:  6.7,
						},
						{
							Date: 1741031636,
							Kwh:  2.4,
						},
						{
							Date: 1741035236,
							Kwh:  7.9,
						},
						{
							Date: 1741038836,
							Kwh:  4.5,
						},
						{
							Date: 1741042436,
							Kwh:  9.3,
						},
						{
							Date: 1741046036,
							Kwh:  1.2,
						},
						{
							Date: 1741049636,
							Kwh:  5.5,
						},
						{
							Date: 1741053236,
							Kwh:  8.4,
						},
						{
							Date: 1741056836,
							Kwh:  3.7,
						},
						{
							Date: 1741060436,
							Kwh:  6.1,
						},
						{
							Date: 1741064036,
							Kwh:  2.8,
						},
						{
							Date: 1741067636,
							Kwh:  7.3,
						},
						{
							Date: 1741071236,
							Kwh:  4.9,
						},
						{
							Date: 1741074836,
							Kwh:  9.6,
						},
						{
							Date: 1741078436,
							Kwh:  1.4,
						},
					},
					Alerts: users.Alert{
						Current: users.CurrentAlert{
							Outage: users.AlertStatus{
								Active: false,
							},
							Tamper: users.AlertStatus{
								Active: false,
							},
						},
						History: []users.AlertEvent{
							{
								Type:      "outage",
								StartDate: time.Now().Unix(),
								EndDate:   time.Now().Unix(),
								Resolved:  true,
							},
						},
					},
					Commands: users.Commands{
						Active: []users.ActiveCommand{
							{
								CommandID:  "adhfoaidfaadf",
								Type:       users.CommandTypeMeterRead,
								IssuedAt:   time.Now().Unix(),
								Parameters: map[string]interface{}{},
								Status:     users.CommandStatusPending,
							},
						},
						History: []users.HistoryCommand{
							{
								CompletedAt: time.Now().Unix(),
								Response:    "Ok",
							},
						},
					},
					Status: users.MeterStatus{
						LastSeen:       time.Now().Unix(),
						GridConnection: true,
						BatteryLevel:   1.0,
					},
				},
			},
			Billing: users.Billing{
				CurrentBill: users.CurrentBill{
					BillingPeriod: users.DateRange{
						Start: time.Now().Unix(),
						End:   time.Now().Unix(),
					},
				},
				PaymentHistory: []users.Payment{
					{
						BillingPeriod: users.DateRange{
							Start: time.Now().Unix(),
							End:   time.Now().Unix(),
						},
						AmountPaid:    decimal.NewFromFloat(3545.05),
						PaymentDate:   time.Now().Unix(),
						PaymentMethod: users.PaymentMethodBankTransfer,
						TransactionID: "edafhkodouiafhoafh",
					},
				},
			},
			Ledger: users.Ledger{
				CurrentBalance: users.Balance{
					Amount:    decimal.NewFromFloat(4000.00),
					UpdatedAt: time.Now().Unix(),
				},
				UpcomingBill: users.UpcomingBill{
					EstimatedAmount: decimal.NewFromFloat(2000.00),
					ProjectionDate:  time.Now().Unix(),
				},
			},

			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		},
		Status: models.AccountStatus{
			IsActive: true,
		},
	}

	// Convert the struct to BSON
	bsonBytes, err := bson.Marshal(acc)
	if err != nil {
		fmt.Println(err)
	}

	// Insert a document
	insertResult, err := controller.Create(context.Background(), bsonBytes)
	if err != nil {
		log.Fatalf("Failed to insert document: %v", err)
	}
	fmt.Printf("Inserted document with ID: %v\n", insertResult.InsertedID)

	//jsonF, _ := json.MarshalIndent(acc, "", "	")
	//fmt.Println(string(jsonF))
}
