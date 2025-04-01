package service

import (
	"fmt"
	"os"
	"strings"
	"uploader/core/upload"
	"uploader/global"
	"uploader/model"
	"uploader/utils"

	"github.com/nguyenthenguyen/docx"
)

type ContractService struct{}

func (*ContractService) CovertContract(data model.ContractModel) (*upload.UploadResult, error) {
	filePath, _ := os.Getwd()
	r, err := docx.ReadDocxFile(filePath + "/file/template/" + data.TemplateName)
	if err != nil {
		return nil, err
	}

	// 下载甲乙双方签名图片
	partyASignaturePath, err := utils.FileHelper.DownloadFile(data.PartyASignature)
	if err != nil {
		return nil, err
	}
	partyBSignaturePath, err := utils.FileHelper.DownloadFile(data.PartyBSignature)
	if err != nil {
		return nil, err
	}

	docx := r.Editable()
	// 姓名
	docx.Replace("partyAName", data.PartyAName, -1)
	docx.Replace("partyBName", data.PartyBName, -1)
	// 价格
	docx.Replace("amount", fmt.Sprintf("%.2f", data.Amount), -1)
	docx.Replace("upperAmount", data.UpperAmount, -1)
	// 定金
	docx.Replace("deposit", fmt.Sprintf("%.2f", data.Deposit), -1)
	docx.Replace("upperDeposit", data.UpperDeposit, -1)
	docx.Replace("payDepositWithinDays", fmt.Sprintf("%d", data.PayDepositWithinDays), -1)
	// 尾款
	docx.Replace("balance", fmt.Sprintf("%.2f", data.Balance), -1)
	docx.Replace("upperBalance", data.UpperBalance, -1)
	docx.Replace("payBalanceWithinDays", fmt.Sprintf("%d", data.PayBalanceWithinDays), -1)
	// 延期交付违约金
	if data.DelayedDeliveryPenalty != nil {
		docx.Replace("delayedDeliveryPenalty", fmt.Sprintf("%.2f", *data.DelayedDeliveryPenalty), -1)
	} else {
		docx.Replace("delayedDeliveryPenalty", "______", -1)
	}
	if data.UpperDelayedDeliveryPenalty != nil {
		docx.Replace("upperDelayedDeliveryPenalty", *data.UpperDelayedDeliveryPenalty, -1)
	} else {
		docx.Replace("upperDelayedDeliveryPenalty", "______", -1)
	}
	// 延期支付违约金
	if data.DelayedPaymentPenalty != nil {
		docx.Replace("delayedPaymentPenalty", fmt.Sprintf("%.2f", *data.DelayedPaymentPenalty), -1)
	} else {
		docx.Replace("delayedPaymentPenalty", "______", -1)
	}
	if data.UpperDelayedPaymentPenalty != nil {
		docx.Replace("upperDelayedPaymentPenalty", *data.UpperDelayedPaymentPenalty, -1)
	} else {
		docx.Replace("upperDelayedPaymentPenalty", "______", -1)
	}
	// 截稿日
	t := strings.Split(data.ExpiresDate, "-")
	docx.Replace("expireDateYear", t[0], -1)
	docx.Replace("expireDateMonth", t[1], -1)
	docx.Replace("expireDateDay", t[2], -1)

	// 签名
	docx.ReplaceImage("word/media/image1.png", partyASignaturePath)
	docx.Replace("partyASignDate", data.PartyASignDate, -1)
	docx.ReplaceImage("word/media/image2.png", partyBSignaturePath)
	docx.Replace("partyBSignDate", data.PartyBSignDate, -1)
	generateFilePath := fmt.Sprintf("%s/%d_%d_contract.docx", global.ServerConfig.System.TemporaryPath,
		data.ID, data.TaskID)
	docx.WriteToFile(generateFilePath)
	r.Close()

	file, err := utils.FileHelper.GetFileHeader(generateFilePath)
	if err != nil {
		return nil, err
	}

	// 上传文件
	uploadResult, err := upload.UploadService.UploadFormFile(upload.UploadFormData{
		Type:    "task",
		ID:      data.TaskID,
		Renamed: true,
		NoGreen: true,
		File:    file,
	})

	if err != nil {
		return nil, err
	}

	// 删除临时文件
	utils.FileHelper.DeleteFile(generateFilePath)
	utils.FileHelper.DeleteFile(partyASignaturePath)
	utils.FileHelper.DeleteFile(partyBSignaturePath)

	return &uploadResult, nil
}
