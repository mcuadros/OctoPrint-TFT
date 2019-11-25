## This project has been archived, but @Z-Bolt has nicely reborn this project at  https://github.com/Z-Bolt/OctoScreen; take a look at it! Looks promising. 

OctoPrint-TFT [![GitHub release](https://img.shields.io/github/release/mcuadros/OctoPrint-TFT.svg)](https://github.com/mcuadros/OctoPrint-TFT/releases) [![license](https://img.shields.io/github/license/mcuadros/OctoPrint-TFT.svg)]()
=============

_OctoPrint-TFT_, a touch interface for TFT touch modules based on GTK+3.

Is a _X application_ to be executed directly in the X Server without any windows
manager, as _frontend of a [OctoPrint](http://octoprint.org) server_ in a Raspberry Pi
equipped with any [TFT Touch module](https://www.waveshare.com/wiki/3.5inch_RPi_LCD_(A)).

Allows you to control your 3D Printer, like you can do with any [TFT/LCD panel](http://reprap.org/wiki/RepRapTouch), but using _OctoPrint_ and a Raspberry Pi.

<img width="480" src="https://user-images.githubusercontent.com/1573114/33559609-a73a969e-d90d-11e7-9cf2-cf212412aaa5.png" />

### These are some of the functionalities supported:

- Print jobs monitoring.
- Temperature and Filament management.
- Jogging operations.

### How this is different from TouchUI?

[TouchUI](http://plugins.octoprint.org/plugins/touchui/), is an amazing plugin
for Octoprint, was created as a responsive design for access to OctoPrint,
from low resolution devices, such as smartphones, tablets, etc.

Executing TouchUI under a RPi w/TFT modules, presents two big problems,
first isn't optimized to be used with resistive touch screens with low resolutions
like 480x320 and second requires a browser to be access, consuming a lot of
resources.

This is the main reason because I develop this X application to be executed
in my 3d printer.

## This version brings some major improvements and fixes:

### Improvements:
- Built-in DPMS control to disable screen blanking
- Status screen optimized, with larger and more visible progress bar and larger buttons
- New button has been added to the temperature controls, allowing to heat only the nozzle up, for example for filament changes
- Redesigned reboot / shutdown buttons in the system menu
- New confirmation dialog when pressing Stop print button

### Fixes:
- Mid-print freezing of the interface has been fixed
- Print/pause/stop buttons now work correctly 
- Random appearance of splash screen after pressing pause / resume has been fixed

Installation
------------

## Install the GUI context
First we need to make sure the GUI context is set up correctly and running on your TFT screen. For this, please follow the instructions, based on your specific installation:


### For Raspbian installation
If you first installed Raspbian, and then manually installed Octoprint:
```sh
sudo apt install lightdm
```
If you have previously attempted to install Octoprint-TFT, during the follwing GUI installation process you could be asked whether to use Lightdm or Octoprint-TFT as your default window manager. Please select Lightdm for now.
After installation is completed, reboot and make sure the GUI works on your TFT screen.
NOTE: if the screen remains blank, you should try reinstalling the TFT screen drivers (depending on your screen make and model).

### For Octopi installation
If you have installed Octopi directly. First you need to install GUI context:
```sh
sudo /home/pi/scripts/install-desktop
```
If you have previously attempted to install Octoprint-TFT, during the follwing GUI installation process you could be asked whether to use Lightdm or Octoprint-TFT as your default window manager. Please select Lightdm for now.
Answer 'yes' to all questions. After installation is completed, reboot and make sure the GUI works on your TFT screen.
NOTE: if the screen remains blank, you should try reinstalling the TFT screen drivers (depending on your screen make and model).

## Uninstall Lightdm window manager
Prior to Octoprint-TFT installation, we need to remove Lightdm window manager, as it could interefere with the successful installation of OctoPrint TFT. In order to do this without removing dependencies that are also required by Octoprint-TFT, run:

```sh
sudo dpkg -r --force-depends lightdm
```

## Dependencies

*OctoPrint-TFT* is based on [Golang](golang.org), usually this means that is
dependency-less, but in this case [GTK+3](https://developer.gnome.org/gtk3/3.0/gtk.html)
is used, this means that GTK+3 libraries are required to be installed on
the system.

If you are using `Raspbian` or any other `Debian` based distribution, GTK+3 can
be installed using:

```sh
sudo apt-get install libgtk-3-0
```
OctoPi does not come with graphical environment, additionally install:

```sh
sudo apt-get install xserver-xorg xinit
```
IMPORTANT!!! In order for the DPMS management to work correctly, you need to install:

```sh
sudo apt-get install x11-xserver-utils
```

## Installation using the Debian installer for Raspbian/OctoPi (recommended)

The recommended way to install *OctoPrint-TFT* is use the `.deb` packages
from the [Releases](https://github.com/darksid3r/OctoPrint-TFT/releases) page. The packages
are available for Debian based distributions such as Raspbian and OctoPi for
versions `Jessie` or `Stretch`.

In order to check which Debian version you have installed, run the following command:

```sh
cat /etc/os-release | grep PRETTY_NAME
```

For example, for a Raspbian Stretch, version 1.1:

```sh
> wget https://github.com/darksid3r/OctoPrint-TFT/releases/download/1.1/octoprint-tft_stretch_1.1.git91fa718-1_armhf.deb 
> sudo dpkg -i octoprint-tft_stretch_1.1.git91fa718-1_armhf.deb 
```
### Please note that in order to get the latest version of Octoprint-TFT for your specific Debian release, go to the "Releases" section of this page.


## Install from source

The compilation and packaging tasks are managed by the [`Makefile`](Makefile)
and backed on [Docker](Dockerfile). Docker is used to avoid installing any other
dependencies since all the operations are done inside of the container.

If you need to install docker inside `Raspbian` or any other linux distrubution
just run:

```sh
curl -fsSL get.docker.com -o get-docker.sh
sh get-docker.sh
```

> You can read more about this at [`docker-install`](https://github.com/docker/docker-install)

To compile the project, assuming that you already cloned this repository, just
execute the `build` target, this will generate in `build` folder all the binaries
and debian packages:

```sh
> make build
> ls -1 build/
```

If you are using `Raspbian` you can install any of the `.deb` generated packages.
If not, just use the compiled binary.

Configuration
-------------

### Basic Configuration

The basic configuration is handled via environment variables, if you are using
the `.deb` package you can configure it at `/etc/octoprint-tft-environment`.

- `OCTOPRINT_CONFIG_FILE` - Location of the OctoPrint's config.yaml file. If empty the file will be searched at the `pi` home folder or the current user. Only used for locally installed OctoPrint servers.

- `OCTOPRINT_HOST` - OctoPrint HTTP address, example `http://localhost:5000`, if OctoPrint is locally installed will be read from the config file.

- `OCTOPRINT_APIKEY` - OctoPrint-TFT expects an [API key]( http://docs.octoprint.org/en/master/api/general.html) to be supplied. This API key can be either the globally configured one or a user specific one if “Access Control”. if OctoPrint is locally installed will be read from the config file.

- `OCTOPRINT_TFT_STYLE_PATH` - Several themes are supported, and style configurations can be done through CSS. This variable defines the location of the application theme.

- `OCTOPRINT_TFT_RESOLUTION` -  Resolution of the application, should be configured to the resolution of your screen, for example `800x480`. By default `480x320`.


### Custom controls and commands

Custom [controls](http://docs.octoprint.org/en/master/configuration/config_yaml.html#controls) to execute GCODE instructions and [commands](http://docs.octoprint.org/en/master/configuration/config_yaml.html#system) to execute shell commands can be defined in the `config.yaml` file.

The controls are limit to static controls without `inputs`.

License
-------

GNU Affero General Public License v3.0, see [LICENSE](LICENSE)

The artwork being use in the at the [default style](`styles/default`) created by [@majurca](https://github.com/majurca) is under the lincese [Attribution 3.0 Unported (CC BY 3.0)](https://creativecommons.org/licenses/by/3.0/)
