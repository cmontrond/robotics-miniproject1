package main

import (
	"fmt"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	g "gobot.io/x/gobot/platforms/dexter/gopigo3"
	"gobot.io/x/gobot/platforms/raspi"
	"time"
)

const (
	LIGHT_IN_REACH  = 1000
	LIGHT_TOO_CLOSE = 3090
	SPEED           = 120
	DIFFERENCE      = 150
)

func stop(gopigo3 *g.Driver) {
	err := gopigo3.SetMotorDps(g.MOTOR_LEFT+g.MOTOR_RIGHT, 0)
	if err != nil {
		fmt.Errorf("Error stopping the robot %+v", err)
	}
}

func left(gopigo3 *g.Driver) {
	err := gopigo3.SetMotorDps(g.MOTOR_LEFT, 0)
	if err != nil {
		fmt.Errorf("Error turning left %+v", err)
	}
	err = gopigo3.SetMotorDps(g.MOTOR_RIGHT, SPEED)
	if err != nil {
		fmt.Errorf("Error turning left %+v", err)
	}
}

func right(gopigo3 *g.Driver) {
	err := gopigo3.SetMotorDps(g.MOTOR_LEFT, SPEED)
	if err != nil {
		fmt.Errorf("Error turning right %+v", err)
	}
	err = gopigo3.SetMotorDps(g.MOTOR_RIGHT, 0)
	if err != nil {
		fmt.Errorf("Error turning right %+v", err)
	}
}

func forward(gopigo3 *g.Driver) {
	err := gopigo3.SetMotorDps(g.MOTOR_LEFT+g.MOTOR_RIGHT, -SPEED)
	if err != nil {
		fmt.Errorf("Error moving forward %+v", err)
	}
}

func blinkLED(gopigo3 *g.Driver) {
	err := gopigo3.SetLED(g.LED_EYE_RIGHT, 0x00, 0x00, 0xFF)
	if err != nil {
		fmt.Println(err)
	}

	time.Sleep(time.Second)

	err = gopigo3.SetLED(g.LED_EYE_RIGHT, 0x00, 0x00, 0x00)
	if err != nil {
		fmt.Println(err)
	}

	time.Sleep(time.Second)
}

func robotRunLoop(gopigo3 *g.Driver, leftLightSensor *aio.GroveLightSensorDriver, rightLightSensor *aio.GroveLightSensorDriver) {

	robotStopped := false

	for {

		// Read value from the left light sensor
		leftLightSensorVal, err := leftLightSensor.Read()

		if err != nil {
			fmt.Errorf("Error reading sensor %+v", err)
		}

		// read value from the right light sensor
		rightLightSensorVal, err := rightLightSensor.Read()

		if err != nil {
			fmt.Errorf("Error reading sensor %+v", err)
		}

		fmt.Println("Right Light Value is ", rightLightSensorVal)
		fmt.Println("Left Light Value is ", leftLightSensorVal)

		// Stop the Robot if too close to the light
		if (rightLightSensorVal >= LIGHT_TOO_CLOSE) && (leftLightSensorVal >= LIGHT_TOO_CLOSE) {
			stop(gopigo3)
			robotStopped = true
			blinkLED(gopigo3)
		}

		if robotStopped == false && leftLightSensorVal > 0 && rightLightSensorVal > 0 {

			rightLeftDifference := rightLightSensorVal - leftLightSensorVal
			leftRightDifference := leftLightSensorVal - rightLightSensorVal

			// If the light comes from the right, turn right and move forward
			if (rightLightSensorVal > leftLightSensorVal) && (rightLightSensorVal >= LIGHT_IN_REACH) && (rightLeftDifference >= DIFFERENCE) {

				right(gopigo3)
				time.Sleep(time.Second)
				stop(gopigo3)
				forward(gopigo3)
				time.Sleep(time.Second)

				// If the light comes from the left, turn left and move forward
			} else if (leftLightSensorVal > rightLightSensorVal) && (leftLightSensorVal >= LIGHT_IN_REACH) && (leftRightDifference >= DIFFERENCE) {

				left(gopigo3)
				time.Sleep(time.Second)
				stop(gopigo3)
				forward(gopigo3)
				time.Sleep(time.Second)

			} else if (rightLightSensorVal >= LIGHT_IN_REACH) && (leftLightSensorVal >= LIGHT_IN_REACH) {

				forward(gopigo3)
				time.Sleep(time.Second)

			} else {
				stop(gopigo3)
			}
		}
	}
}

func main() {

	raspiAdaptor := raspi.NewAdaptor()
	gopigo3 := g.NewDriver(raspiAdaptor)

	leftLightSensor := aio.NewGroveLightSensorDriver(gopigo3, "AD_2_1")
	rightLightSensor := aio.NewGroveLightSensorDriver(gopigo3, "AD_1_1")

	mainRobotFunc := func() {
		robotRunLoop(gopigo3, leftLightSensor, rightLightSensor)
	}

	robot := gobot.NewRobot("miniProject",
		[]gobot.Connection{raspiAdaptor},
		[]gobot.Device{gopigo3, leftLightSensor, rightLightSensor},
		mainRobotFunc,
	)

	err := robot.Start()

	if err != nil {
		fmt.Errorf("Error starting the Robot %+v", err)
	}
}
