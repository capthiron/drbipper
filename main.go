package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	bip39 "github.com/tyler-smith/go-bip39/wordlists"
)

type mode int

const (
	modeMenu mode = iota
	modePassword
	modeMnemonic
	modeResult
)

type model struct {
	mode     mode
	choice   string
	password string
	mnemonic []string
	cursor   int
	result   []string
	err      string
}

func initialModel() *model {
	return &model{
		mode:     modeMenu,
		mnemonic: make([]string, 24),
		cursor:   0,
	}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.mode {
		case modeMenu:
			switch msg.String() {
			case "1":
				m.choice = "encrypt"
				m.mode = modePassword
			case "2":
				m.choice = "decrypt"
				m.mode = modePassword
			case "q":
				return m, tea.Quit
			}
		case modePassword:
			switch msg.Type {
			case tea.KeyEnter:
				m.mode = modeMnemonic
			case tea.KeyBackspace:
				if len(m.password) > 0 {
					m.password = m.password[:len(m.password)-1]
				}
			case tea.KeyRunes:
				m.password += string(msg.Runes)
			default:
				panic("unhandled default case")
			}
		case modeMnemonic:
			switch msg.Type {
			case tea.KeyEnter:
				if m.cursor < 24 {
					if indexOf(m.mnemonic[m.cursor], bip39.English) == -1 {
						// Word not found in the wordlist, request again
						m.mnemonic[m.cursor] = ""
						m.err = "Invalid word, please try again."
					} else {
						m.cursor++
						m.err = ""
					}
				}
				if m.cursor == 24 {
					m.mode = modeResult
					m.processMnemonic()
				}
			case tea.KeyBackspace:
				if m.mnemonic[m.cursor] != "" {
					m.mnemonic[m.cursor] = ""
					break
				}
				if m.cursor > 0 {
					m.cursor--
				}
			case tea.KeyRunes:
				if m.cursor < 24 {
					m.mnemonic[m.cursor] += string(msg.Runes)
				}
			default:
				panic("unhandled default case")
			}
		case modeResult:
			if msg.Type == tea.KeyEnter {
				return m, tea.Quit
			}
		default:
			panic("unhandled default case")
		}
	}
	return m, nil
}

func (m *model) View() string {
	switch m.mode {
	case modeMenu:
		return "Choose an option:\n1. Encrypt mnemonic\n2. Decrypt mnemonic\nq. Quit\n"
	case modePassword:
		return fmt.Sprintf("Enter password: %s", strings.Repeat("*", len(m.password)))
	case modeMnemonic:
		errorMsg := ""
		if m.err != "" {
			errorMsg = fmt.Sprintf("\n\nError: %s", m.err)
		}
		enteredWords := strings.Join(m.mnemonic[:m.cursor], " ")
		return fmt.Sprintf("Enter word %d/24: %s\n%s%s", m.cursor+1, m.mnemonic[m.cursor], enteredWords, errorMsg)
	case modeResult:
		return fmt.Sprintf("Result:\n%s\n\nPress Enter to exit and clear the terminal.", strings.Join(m.result, " "))
	default:
		return ""
	}
}

func (m *model) processMnemonic() {
	hash := hashPassword(m.password)
	shifts := generateShiftValues(hash)

	if m.choice == "encrypt" {
		m.result = Encrypt(m.mnemonic, shifts)
	} else if m.choice == "decrypt" {
		m.result = Decrypt(m.mnemonic, shifts)
	}
}

func clearTerminal() error {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func hashPassword(password string) []byte {
	hash := sha256.Sum256([]byte(password))
	return hash[:]
}

func generateShiftValues(hash []byte) []int {
	shiftValues := make([]int, 24)
	for i := 0; i < 24; i++ {
		start := (i * 4) % len(hash)
		end := start + 4
		shiftValues[i] = int(binary.BigEndian.Uint32(hash[start:end]) % 2048)
	}
	return shiftValues
}

func Encrypt(mnemonic []string, shifts []int) []string {
	shiftedMnemonic := make([]string, len(mnemonic))
	for i, word := range mnemonic {
		index := indexOf(word, bip39.English)
		if index == -1 {
			fmt.Printf("Word %s not found in BIP-39 wordlist\n", word)
			os.Exit(1)
		}
		newIndex := (index + shifts[i]) % len(bip39.English)
		shiftedMnemonic[i] = bip39.English[newIndex]
	}
	return shiftedMnemonic
}

func Decrypt(shiftedMnemonic []string, shifts []int) []string {
	originalMnemonic := make([]string, len(shiftedMnemonic))
	for i, word := range shiftedMnemonic {
		index := indexOf(word, bip39.English)
		if index == -1 {
			fmt.Printf("Word %s not found in BIP-39 wordlist\n", word)
			os.Exit(1)
		}
		newIndex := (index - shifts[i] + len(bip39.English)) % len(bip39.English)
		originalMnemonic[i] = bip39.English[newIndex]
	}
	return originalMnemonic
}

func indexOf(word string, wordlist []string) int {
	for i, w := range wordlist {
		if w == word {
			return i
		}
	}
	return -1
}

func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	err := clearTerminal()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(2)
	}
}
