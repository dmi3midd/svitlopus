package dockerutil

import (
	"fmt"
	"log"
	"os/exec"
	"svitlopus/internal/config"
)

type DockerUtil interface {
	IsDockerInstalled() bool
	IsImagePulled() bool
	PullImage()
	DoesContainerExist() bool
	RunDockerContainer() error
	RunDockerPipeline() error
}

type dockerUtil struct {
	cfg *config.Docker
}

func NewDockerUtil(cfg *config.Docker) DockerUtil {
	return &dockerUtil{
		cfg: cfg,
	}
}

func (u *dockerUtil) IsDockerInstalled() bool {
	cmd := exec.Command("docker", "-v")
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func (u *dockerUtil) IsImagePulled() bool {
	image := u.cfg.Image
	cmd := exec.Command("docker", "image", "inspect", image)
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func (u *dockerUtil) PullImage() {
	image := u.cfg.Image
	cmd := exec.Command("docker", "pull", image)
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("failed to pull image: %v\n", err)
		return
	}
}

func (u *dockerUtil) DoesContainerExist() bool {
	containerName := u.cfg.ContainerName
	cmd := exec.Command("docker", "inspect", containerName)
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	return cmd.Run() == nil
}

func (u *dockerUtil) RunDockerContainer() error {
	port := u.cfg.Port
	containerName := u.cfg.ContainerName
	volume := u.cfg.Volume
	apiId := u.cfg.ApiId
	apiHash := u.cfg.ApiHash
	image := u.cfg.Image
	restartAlways := u.cfg.RestartAlways
	restartFlag := "always"
	if !restartAlways {
		restartFlag = "no"
	}
	args := []string{
		"run", "-d",
		"-p", fmt.Sprintf("%d:%d", port, port),
		"--name", containerName,
		"--restart", restartFlag,
		"-v", fmt.Sprintf("%s:/var/lib/telegram-bot-api", volume),
		"-e", fmt.Sprintf("TELEGRAM_API_ID=%s", apiId),
		"-e", fmt.Sprintf("TELEGRAM_API_HASH=%s", apiHash),
		image,
	}

	cmd := exec.Command("docker", args...)
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run container: %v", err)
	}
	return nil
}

func (u *dockerUtil) RunDockerPipeline() error {
	cond := u.IsDockerInstalled()
	if !cond {
		return fmt.Errorf("docker is not installed")
	}
	cond = u.IsImagePulled()
	if !cond {
		u.PullImage()
	}
	cond = u.DoesContainerExist()
	if !cond {
		u.RunDockerContainer()
	}
	return nil
}
