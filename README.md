# STT

STT is a command line application which wraps [AssemblyAI](https://www.assemblyai.com/)
and performs the necessary requests to convert a voice recording to structured data.

## Installation

1. [Create an account with AssemblyAI](https://app.assemblyai.com/) and get your API key from the developer dashboard under the "Account" section.
1. Make sure [Go](https://go.dev/dl/) is installed on your system and included in your `PATH`.
1. Clone this repository into a directory on your local machine and `cd` into it.
1. Create a `.env` in the root of your project and add the following text: `ASSEMBLY_AI_KEY=YOUR_API_KEY`
1. Compile the program by running: `go build .`.

To make the package globally executable run:

`$ go install .`

## Usage

Once the app is compiled, you can run it like so:

`$ ./stt [path_to_audio_file].m4a`

This will create a file called `out.json` in the local directory containing all values from AssemblyAI in a JSON object.
