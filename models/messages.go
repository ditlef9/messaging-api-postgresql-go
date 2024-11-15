// File: models/messages.go

package models

import (
	"database/sql"
	"ekeberg.com/messaging-api-postgresql-go/db"
	"strconv"
)

type Message struct {
	Id                int          `json:"msg_id"`
	Platform          string       `json:"msg_platform"`
	ExternalId        string       `json:"msg_external_id"`
	CreatedAt         string       `json:"msg_created_at"`
	Language          string       `json:"msg_language"`
	URL               string       `json:"msg_url"`
	Content           string       `json:"msg_content"`
	ExternalAccountId string       `json:"msg_external_account_id"`
	Attachments       []Attachment `json:"attachments"`
}
type Attachment struct {
	AttachmentId              int    `json:"attachment_id"`
	AttachmentMsgId           int    `json:"attachment_msg_id"`
	AttachmentExternalId      string `json:"attachment_external_id"`
	AttachmentURL             string `json:"attachment_url"`
	AttachmentType            string `json:"attachment_type"`
	AttachmentMetaDescription string `json:"attachment_meta_description"`
}

func GetMessages(count int) ([]Message, error) {
	// Initialize the base query for messages
	messageQuery := "SELECT msg_id, msg_platform, msg_external_id, msg_created_at, msg_language, msg_url, msg_content, msg_external_account_id FROM messages_index ORDER BY msg_created_at DESC"

	// Modify the query to include a LIMIT if count is greater than zero
	if count > 0 {
		messageQuery += " LIMIT " + strconv.Itoa(count)
	} else if count == 0 {
		// If count is 0, default it to 10000
		count = 10000
		messageQuery += " LIMIT " + strconv.Itoa(count)
	}

	// Execute the query
	messageRows, err := db.DB.Query(messageQuery)
	if err != nil {
		return nil, err
	}
	defer messageRows.Close()

	// Initialize an empty slice to store messages
	messages := make([]Message, 0)

	// Iterate through the message rows and fetch attachments
	for messageRows.Next() {
		sqlMessage := Message{}
		err = messageRows.Scan(&sqlMessage.Id, &sqlMessage.Platform, &sqlMessage.ExternalId, &sqlMessage.CreatedAt, &sqlMessage.Language, &sqlMessage.URL, &sqlMessage.Content, &sqlMessage.ExternalAccountId)
		if err != nil {
			return nil, err
		}

		// Fetch attachments for the current message
		attachmentQuery := "SELECT attachment_id, attachment_msg_id, attachment_external_id, attachment_url, attachment_type, attachment_meta_description FROM messages_attachments WHERE attachment_msg_id = $1"
		attachmentRows, err := db.DB.Query(attachmentQuery, sqlMessage.Id)
		if err != nil {
			return nil, err
		}
		defer attachmentRows.Close()

		attachments := make([]Attachment, 0)
		for attachmentRows.Next() {
			attachment := Attachment{}
			err = attachmentRows.Scan(&attachment.AttachmentId, &attachment.AttachmentMsgId, &attachment.AttachmentExternalId, &attachment.AttachmentURL, &attachment.AttachmentType, &attachment.AttachmentMetaDescription)
			if err != nil {
				return nil, err
			}
			attachments = append(attachments, attachment)
		}
		sqlMessage.Attachments = attachments

		messages = append(messages, sqlMessage)
	}

	// Check for any errors encountered during row iteration
	if err = messageRows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func GetMessageById(id string) (Message, error) {
	// Prepare the query for the message
	stmt, err := db.DB.Prepare("SELECT msg_id, msg_platform, msg_external_id, msg_created_at, msg_language, msg_url, msg_content, msg_external_account_id FROM messages_index WHERE msg_id = $1")
	if err != nil {
		return Message{}, err
	}

	sqlMessage := Message{}
	sqlErr := stmt.QueryRow(id).Scan(&sqlMessage.Id, &sqlMessage.Platform, &sqlMessage.ExternalId, &sqlMessage.CreatedAt, &sqlMessage.Language, &sqlMessage.URL, &sqlMessage.Content, &sqlMessage.ExternalAccountId)
	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return Message{}, nil
		}
		return Message{}, sqlErr
	}

	// Fetch attachments for the message
	attachmentQuery := "SELECT attachment_id, attachment_msg_id, attachment_external_id, attachment_url, attachment_type, attachment_meta_description FROM messages_attachments WHERE attachment_msg_id = $1"
	attachmentRows, err := db.DB.Query(attachmentQuery, sqlMessage.Id)
	if err != nil {
		return Message{}, err
	}
	defer attachmentRows.Close()

	attachments := make([]Attachment, 0)
	for attachmentRows.Next() {
		attachment := Attachment{}
		err = attachmentRows.Scan(&attachment.AttachmentId, &attachment.AttachmentMsgId, &attachment.AttachmentExternalId, &attachment.AttachmentURL, &attachment.AttachmentType, &attachment.AttachmentMetaDescription)
		if err != nil {
			return Message{}, err
		}
		attachments = append(attachments, attachment)
	}
	sqlMessage.Attachments = attachments

	return sqlMessage, nil
}
