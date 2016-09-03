// Pin Functions
#define FORWARD_PIN   (9)
#define BACKWARD_PIN (10)
#define LEFT_PIN     (11)
#define RIGHT_PIN    (12)

// Bits to indicate FORWARD, BACKWARD, LEFT, and RIGHT
#define FORARD_BIT   (1) // b'0001' (binary)
#define BACKWARD_BIT (2) // b'0010'
#define LEFT_BIT     (4) // b'0100'
#define RIGHT_BIT    (8) // b'1000'

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
void driveCar(struct Command &newCmd)
{
    // If forward and backward are both enabled, error, remove the backward bit set
    if ((newCmd.direction & FORARD_BIT) && (newCmd.direction & BACKWARD_BIT)) {
        newCmd.direction -= BACKWARD_BIT;
    }
    
    // If left and right are both enabled, error, remove the right bit set
    if ((newCmd.direction & LEFT_BIT) && (newCmd.direction & RIGHT_BIT)) {
        newCmd.direction -= RIGHT_BIT;
    }
    
    // Drive forward if enabled
    if (newCmd.direction & FORARD_BIT) {
        // Note: newCmd.data2 is the speed, a PWM value specified in range 0 - 255, 255 = MAX
        analogWrite(FORWARD_PIN, newCmd.speed);
    } else {
        analogWrite(FORWARD_PIN, 0);
    }
    
    // Drive backward if enabled
    if (newCmd.direction & BACKWARD_BIT) {
        analogWrite(BACKWARD_PIN, newCmd.speed);
    } else {
        analogWrite(BACKWARD_PIN, 0);
    }
    
    // Drive left if enabled
    if (newCmd.direction & LEFT_BIT) {
        digitalWrite(LEFT_PIN, HIGH);
    } else {
        digitalWrite(LEFT_PIN, LOW);
    }
    
    // Drive right if enabled
    if (newCmd.direction & RIGHT_BIT) {
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
    } else {
      Serial.print("Got number of bytes: ");
      Serial.println(in, DEC);
    }
  }
}
