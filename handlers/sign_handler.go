package handler

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"net/http"
	"os"
	"path/filepath"

	"go_digital_sign/model"

	"github.com/nfnt/resize"
	"github.com/skip2/go-qrcode"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var signCollection *mongo.Collection

func SetSignCollection(c *mongo.Collection) {
	signCollection = c
}

type SignRequest struct {
	UserUUID      string `json:"user_uuid"`
	TitleDoc      string `json:"title_document"`
	Category      string `json:"category"`
	Description   string `json:"description"`
	SignTimestamp string `json:"sign_timestamp"`
}

type SignResponse struct {
	TitleDoc       string `json:"title_document"`
	Category       string `json:"category"`
	Description    string `json:"description"`
	SignTimestamp  string `json:"sign_timestamp"`
	LinktoVerify   string `json:"link_verify"`
	StorageAddress string `json:"storage_address"`
	Status         string `json:"status"`
}

func GenerateSignature(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

func createQRCodeWithLogo(data, filepath, logoPath string) error {
	qr, err := qrcode.New(data, qrcode.Medium)
	if err != nil {
		return err
	}

	qrImage := qr.Image(256)

	logoFile, err := os.Open(logoPath)
	if err != nil {
		return err
	}
	defer logoFile.Close()

	logoImage, _, err := image.Decode(logoFile)
	if err != nil {
		return err
	}

	logoResized := resize.Resize(64, 64, logoImage, resize.Lanczos3)

	offset := image.Pt((qrImage.Bounds().Dx()-logoResized.Bounds().Dx())/2, (qrImage.Bounds().Dy()-logoResized.Bounds().Dy())/2)

	b := qrImage.Bounds()
	m := image.NewRGBA(b)
	draw.Draw(m, b, qrImage, image.Point{}, draw.Src)
	draw.Draw(m, logoResized.Bounds().Add(offset), logoResized, image.Point{}, draw.Over)

	outfile, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer outfile.Close()

	err = jpeg.Encode(outfile, m, nil)
	if err != nil {
		return err
	}

	return nil
}

func SignDocument(w http.ResponseWriter, r *http.Request) {
	var req SignRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	signatureData := fmt.Sprintf("%s|%s|%s|%s|%s", req.UserUUID, req.TitleDoc, req.Category, req.Description, req.SignTimestamp)
	signature := GenerateSignature(signatureData)

	domain := os.Getenv("HOSTNAME")
	linktoVerify := fmt.Sprintf("%s/api/signatureverify?credentialcode=%s", domain, signature)

	uuid := req.UserUUID // Using UserUUID as UUID for storage
	storagePath := filepath.Join("data", fmt.Sprintf("%s.jpg", uuid))

	err = createQRCodeWithLogo(linktoVerify, storagePath, "assets/img/logo.jpg")
	if err != nil {
		http.Error(w, "Failed to create QR code", http.StatusInternalServerError)
		return
	}

	sign := model.Sign{
		UUID:           uuid,
		UserUUID:       req.UserUUID,
		TitleDoc:       req.TitleDoc,
		Category:       req.Category,
		Description:    req.Description,
		Signature:      signature,
		LinktoVerify:   linktoVerify,
		SignTimestamp:  req.SignTimestamp,
		StorageAddress: storagePath,
	}

	_, err = signCollection.InsertOne(context.TODO(), sign)
	if err != nil {
		http.Error(w, "Failed to save the document", http.StatusInternalServerError)
		return
	}

	response := SignResponse{
		TitleDoc:       req.TitleDoc,
		Category:       req.Category,
		Description:    req.Description,
		SignTimestamp:  req.SignTimestamp,
		LinktoVerify:   linktoVerify,
		StorageAddress: storagePath,
		Status:         "success",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

type VerifyResponse struct {
	Status        string `json:"status"`
	UserUUID      string `json:"user_uuid"`
	TitleDoc      string `json:"title_document"`
	Category      string `json:"category"`
	Description   string `json:"description"`
	SignTimestamp string `json:"sign_timestamp"`
}

func VerifySignature(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	credentialCode := queryParams.Get("credentialcode")

	if credentialCode == "" {
		http.Error(w, "Missing credential code", http.StatusBadRequest)
		return
	}

	var sign model.Sign
	err := signCollection.FindOne(context.TODO(), bson.M{"signature": credentialCode}).Decode(&sign)
	if err != nil {
		response := VerifyResponse{
			Status: "Signature not found or invalid",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := VerifyResponse{
		Status:        "Signature is valid",
		UserUUID:      sign.UserUUID,
		TitleDoc:      sign.TitleDoc,
		Category:      sign.Category,
		Description:   sign.Description,
		SignTimestamp: sign.SignTimestamp,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
