# drbipper

Welcome to **drbipper**! 🎉 This CLI tool is inspired by the board game **Doktor Bibber** and the BIP-39 mnemonic phrases. Just like the game, drbipper requires precision and care, but instead of removing funny bones, you'll be encrypting and decrypting mnemonic phrases! 🧩🔐

## Features

- Encrypt mnemonic phrases using a password. 🔒
- Decrypt mnemonic phrases using a password. 🔓

## Installation

1. Clone the repository:
   ~~~sh
   git clone https://github.com/yourusername/drbipper.git
   cd drbipper
   ~~~

2. Install dependencies:
   ~~~sh
   go mod tidy
   ~~~

3. Build the application:
   ~~~sh
   go build -o drbipper
   ~~~

## Usage

Run the application:
~~~sh
./drbipper
~~~

Follow the on-screen instructions to encrypt or decrypt your mnemonic phrase.

## Example

1. Choose an option:
   ~~~
   Choose an option:
   1. Encrypt mnemonic
   2. Decrypt mnemonic
   q. Quit
   ~~~

2. Enter your password:
   ~~~
   Enter password: ********
   ~~~

3. Enter your mnemonic words one by one. If a word is invalid, you will be prompted to enter it again:
   ~~~
   Enter word 1/24: oppose
   Enter word 2/24: duck
   ...
   ~~~

4. View the result:
   ~~~
   Result:
   <encrypted or decrypted mnemonic>
   
   Press Enter to exit and clear the terminal.
   ~~~

## Testing

Run the tests:
~~~sh
go test ./...
~~~

## Security Notice

**Caution:** Handling crypto passphrases requires careful consideration of security practices. Be aware of the following:

- **Environment:** Ensure that the environment where this code is run is secure. Avoid running this code on shared or untrusted devices.
- **Input Devices:** Be cautious about the devices into which you enter your passphrases. Ensure they are free from malware and keyloggers.
- **Storage:** Do not store your passphrases or passwords in plain text. Use secure storage mechanisms.
- **Code Understanding:** Only run this script if you understand the code and what it does. Anyone can fork it, turn it malicious, and trick you into using it if you don't understand the underlying code. For safety, run this script on an air-gapped computer that is not connected to the internet. A better option is to do it manually with the "mapping_table.txt" file. Write down the encrypted seed words by hand, do not print them. Even better, stamp or engrave them on titanium plates to protect from fire or water damage.
- **Security Level:** Note that the encrypted words/numbers are not cryptographically secure and can be bruteforced. They provide some protection from common thieves and extra time to react in case of theft.

## Disclaimer

**Warning:** There is no guarantee that this code works! Use it at your own risk and only if you fully understand the concepts applied.

## Credits

This project was inspired by [Seedshift](https://github.com/mifunetoshiro/Seedshift), which uses dates as the seed for the encryption (shift operation) instead of passwords.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
