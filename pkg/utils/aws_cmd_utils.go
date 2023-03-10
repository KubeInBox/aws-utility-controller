package utils

import (
	"context"
	"os/exec"
	"strings"
	"time"

	"github.com/go-logr/logr"
)

func runCMD(logger logr.Logger, executeCMD string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.V(1).Info("running command", "cmd", executeCMD)
	cmd := exec.CommandContext(ctx, "sh", "-c", executeCMD)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error(err, "unable to start instances", "stderr", cmd.Stderr)
		return err
	}
	logger.V(2).Info("cmd output", "stdout", string(out))
	return nil
}

func StartEc2Instance(logger logr.Logger, instanceIDs []string) error {
	toBeExecuted := "aws ec2 start-instances --instance-ids " + strings.Join(instanceIDs, ",")
	if err := runCMD(logger, toBeExecuted); err != nil {
		return err
	}
	logger.Info("successfully started ec2 instances")
	return nil
}

func StopEc2Instance(logger logr.Logger, instanceIDs []string) error {
	toBeExecuted := "aws ec2 stop-instances --instance-ids " + strings.Join(instanceIDs, ",")
	if err := runCMD(logger, toBeExecuted); err != nil {
		return err
	}
	logger.Info("successfully stopped ec2 instances")
	return nil
}

// TODO:
// p2: parse aws credentials from end user.
// p2: Validations on CRs Fields.
// p1: start/stop ec2 in counterpart of schedule time window. akshay
// p1: what if user want onDemand in schedule window `? joy
// for every onDemand req check whether corresponding scheduled is available or not.
//if ondemand comes when "Scheduled/InTimeWindow" ||  "Scheduled/OutTimeWindow"
//	- pause scheudled operation.
//    - perform ondemand operation.
// p1: avoid processing on-demand multiple times in case of success. shani
