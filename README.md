# Linebot Server Felice
![Felice](./images/Felice.png)
Linebot Server Felice is a Golang-based Line bot server that allows you to create your own personalized and unique Line bot.

## Table of Contents

- [Introduction](#introduction)
- [Technologies](#technologies)
- [Getting Started](#getting-started)
- [Environment Variables](#environment-variables)
- [Usage Examples](#usage-examples)
- [Contributing](#contributing)
- [License](#license)

## Introduction

Linebot Server Felice is a project aimed at building your own Line bot with a touch of personality. The bot, named "Felice," is an imaginative and adventurous female character known for her passion for exploration and discovery. She has a cheerful personality and boundless energy, but sometimes responds with "Ahahaha" to ease awkward moments when she can't understand others.

## Technologies

This project is built using [go-linebot-service](https://github.com/Zhima-Mochi/go-linebot-service). It also utilizes the [OpenAI API](https://openai.com/) for natural language processing capabilities.

## Getting Started

To get started with Linebot Server Felice, follow these steps:

1. Clone the repository to your local machine.
2. Install the required dependencies using `go mod`.
3. Set up the necessary environment variables (details provided in the [Environment Variables](#environment-variables) section).
4. Deploy the Line bot server using the provided Dockerfiles or your preferred deployment method.
5. Interact with your personalized Line bot!

## Environment Variables

Before running the Linebot Server Felice, you need to set the following environment variables:

- `LINE_CHANNEL_SECRET`: (Your Line channel secret)
- `LINE_CHANNEL_TOKEN`: (Your Line channel access token)
- `OPENAI_API_KEY`: (Your OpenAI API key)
- `CACHE_URL`: (URL for caching purposes)
- `LINEBOT_PORT`: (Port on which the Line bot server will run)
- `LINE_ADMIN_USER_ID_LIST`: (Comma-separated list of Line user IDs for admin access)

Please ensure that you have obtained the necessary API keys and tokens for your Line channel and OpenAI account.

## Usage Examples

Here are some examples of how you can interact with Linebot Server Felice:

- Example 1: [Describe an interesting place you've visited]
- Example 2: [Tell me a joke]
- Example 3: [What are your favorite hobbies?]

Feel free to experiment with different questions and prompts to see how Felice responds!

## Contributing

Contributions to Linebot Server Felice are welcome! If you have any improvements, bug fixes, or new features to add, please submit a pull request. For major changes, please open an issue to discuss the proposed changes.

## License

Linebot Server Felice is open-source and distributed under the [MIT License](LICENSE).
