package utils

import (
	"context"
	"os/exec"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func runCMD(executeCMD string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Infof("running command %s", executeCMD)
	cmd := exec.CommandContext(ctx, "sh", "-c", executeCMD)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.WithError(err).Errorf("unable to start instances, error: %v", cmd.Stderr)
		return err
	}
	log.Debugf("output  output=%s", string(out))
	return nil
}

func StartEc2Instance(instanceIDs []string) error {
	toBeExecuted := "aws ec2 start-instances --instance-ids " + strings.Join(instanceIDs, ",")
	if err := runCMD(toBeExecuted); err != nil {
		return err
	}
	log.Info("successfully started ec2 instances")
	return nil
}

func StopEc2Instance(instanceIDs []string) error {
	toBeExecuted := "aws ec2 stop-instances --instance-ids " + strings.Join(instanceIDs, ",")
	if err := runCMD(toBeExecuted); err != nil {
		return err
	}
	log.Info("successfully stopped ec2 instances")
	return nil
}