# Tomatigo

The *simplest* Pomodoro timer for your terminal.

![Select your pomodoro style](tomatigo-selection.png)
![Focus time](tomatigo-focus.png)

## Installation

Tomatigo uses Oto to play sounds. So check the installation documentation:

- [Oto](https://github.com/ebitengine/oto?tab=readme-ov-file#prerequisite)

### Installing tomatigo

#### Option A (requires Go)
```
git clone <https://github.com/xgrimauu/tomatigo.git>
cd tomatigo
go build .
./tomatigo
```
#### Option B

Downloading binary option: Simply download the binary and run ./tomatigo

### Use Tomatigo

- `<space>` will enter the selection and pause the timer.
- `<q>` will quit return to timer selection, or the will exit the program if pressed in the selection screen

- Tomatigo supports both Vim keybindings (h,j,k,l) and Arrows for navigation.
