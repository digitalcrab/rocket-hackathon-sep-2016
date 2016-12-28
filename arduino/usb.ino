// Pin Functions
#define FORWARD_PIN   (9)
#define BACKWARD_PIN (10)
#define LEFT_PIN     (11)
#define RIGHT_PIN    (12)

// Bits to indicate FORWARD, BACKWARD, LEFT, and RIGHT
#define FORARD_BIT   (1) // '0001'
#define BACKWARD_BIT (2) // '0010'
#define LEFT_BIT     (4) // '0100'
#define RIGHT_BIT    (8) // '1000'

// Each command is 2 bytes in size
struct Command
{
    byte direction;
    byte speed;
};

void setup()
{
    // Setup Pin I/O Functions
    pinMode(FORWARD_PIN, OUTPUT);
    pinMode(BACKWARD_PIN, OUTPUT);
    pinMode(LEFT_PIN, OUTPUT);
    pinMode(RIGHT_PIN, OUTPUT);
    
    // Initialize Serial
    Serial.begin(9600);
}

// Decodes a command struct, does some error checking, and controls the Arduino pins
void driveCar(struct Command &cmd)
{
    // If forward and backward are both enabled, error, remove the backward bit set
    if ((cmd.direction & FORARD_BIT) && (cmd.direction & BACKWARD_BIT)) {
        cmd.direction -= BACKWARD_BIT;
    }
    
    // If left and right are both enabled, error, remove the right bit set
    if ((cmd.direction & LEFT_BIT) && (cmd.direction & RIGHT_BIT)) {
        cmd.direction -= RIGHT_BIT;
    }
    
    // Drive forward if enabled
    if (cmd.direction & FORARD_BIT) {
        // Note: a PWM value specified in range 0 - 255, 255 = MAX
        analogWrite(FORWARD_PIN, cmd.speed);
    } else {
        analogWrite(FORWARD_PIN, 0);
    }
    
    // Drive backward if enabled
    if (cmd.direction & BACKWARD_BIT) {
        analogWrite(BACKWARD_PIN, cmd.speed);
    } else {
        analogWrite(BACKWARD_PIN, 0);
    }
    
    // Drive left if enabled
    if (cmd.direction & LEFT_BIT) {
        digitalWrite(LEFT_PIN, HIGH);
    } else {
        digitalWrite(LEFT_PIN, LOW);
    }
    
    // Drive right if enabled
    if (cmd.direction & RIGHT_BIT) {
        digitalWrite(RIGHT_PIN, HIGH);
    } else {
        digitalWrite(RIGHT_PIN, LOW);
    }
}

byte buf[2];

void loop()
{
  if (Serial.available() >= 2) {
    int in = Serial.readBytes(buf, 2);
    if (in == 2) {
      Serial.println("Got 2 bytes");
      Serial.println(buf[0], DEC);
      Serial.println(buf[1], DEC);
      Command cmd;
      cmd.direction = buf[0];
      cmd.speed = buf[1];
      driveCar(cmd); 
    }
  }
}
