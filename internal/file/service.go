package file

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	telegramadapter "svitlopus/internal/adapters/telegram_adapter"
	"svitlopus/internal/chats"
	"svitlopus/internal/folder"
	"time"

	"github.com/rs/xid"
	"gopkg.in/telebot.v4"
)

var (
	ErrParentNotFound   = errors.New("parent folder not found")
	ErrFileAlreadyExist = errors.New("file with this tutle already exist")
)

type FileService interface {
	GetFile(ctx context.Context, fileId string) (*File, error)
	GetSubfiles(ctx context.Context, fileId string, limit, offset int) ([]File, error)
	RenameFile(ctx context.Context, id, newTitle string) (*File, error)
	MoveFile(ctx context.Context, id, newParentId string) (*File, error)
	DeleteFile(ctx context.Context, id string) error
	UploadFile(ctx context.Context, file multipart.File, title, parentId string) (*File, error)
	DownloadFile(ctx context.Context, id string) *telebot.File
}

type fileService struct {
	telegramUtil     telegramadapter.TelegramClient
	chatManager      chats.ChatManager
	fileRepository   FileRepository
	folderRepository folder.FolderRepository
}

func NewFileService(
	telegramUtil telegramadapter.TelegramClient,
	chatManager chats.ChatManager,
	fileRepository FileRepository,
	folderRepository folder.FolderRepository,
) FileService {
	return &fileService{
		telegramUtil:     telegramUtil,
		chatManager:      chatManager,
		fileRepository:   fileRepository,
		folderRepository: folderRepository,
	}
}

// GetFile implements [FileService].
func (s *fileService) GetFile(ctx context.Context, fileId string) (*File, error) {
	panic("unimplemented")
}

// GetFiles implements [FileService].
func (s *fileService) GetSubfiles(ctx context.Context, fileId string, limit int, offset int) ([]File, error) {
	panic("unimplemented")
}

// RenameFile implements [FileService].
func (s *fileService) RenameFile(ctx context.Context, id string, newTitle string) (*File, error) {
	panic("unimplemented")
}

// MoveFile implements [FileService].
func (s *fileService) MoveFile(ctx context.Context, id string, newParentId string) (*File, error) {
	panic("unimplemented")
}

// DeleteFile implements [FileService].
func (s *fileService) DeleteFile(ctx context.Context, id string) error {
	panic("unimplemented")
}

func (s *fileService) UploadFile(ctx context.Context, file multipart.File, title, parentId string) (*File, error) {
	op := "FileService.UploadFile"

	// Check parent folder
	if _, err := s.folderRepository.GetById(ctx, parentId); err != nil {
		if errors.Is(err, folder.ErrNoFolder) {
			return nil, fmt.Errorf("%s: %w", op, ErrParentNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Check for title collision
	_, err := s.fileRepository.GetByTitleAndParentId(ctx, title, parentId)
	if err == nil {
		return nil, fmt.Errorf("%s: %w", op, ErrFileAlreadyExist)
	}
	if !errors.Is(err, ErrNoFile) {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Create file
	chatId, _ := s.chatManager.NextChat()
	f, err := s.telegramUtil.UploadFile(file, chatId)
	id := xid.New().String()
	newFile := File{
		Id:        id,
		Title:     title,
		Mime:      f.MimeType,
		Size:      f.Size,
		ParentId:  parentId,
		FileId:    f.TgFileId,
		MessageId: f.TgMsgId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return &newFile, nil
}

// DownloadFile implements [FileService].
func (s *fileService) DownloadFile(ctx context.Context, id string) *telebot.File {
	panic("unimplemented")
}
