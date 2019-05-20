#include <plasma/client.h>
#include <iostream>
using namespace plasma;

int main(int argc, char** argv) {
  // Start up and connect a Plasma client.
  PlasmaClient client;
  auto status = client.Connect("/tmp/plasma", "");
  if (!status.ok()) {
      std::cerr << "error: " << status.message() << "\n";
      std::exit(1);
  }
  // Disconnect the Plasma client.
//  ARROW_CHECK_OK(client.Disconnect());
}