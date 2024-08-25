# G3N - Go 3D Game Engine

**G3N** is an OpenGL 3D Game Engine written in Go.
It can be used to write cross-platform Go applications that show rich and dynamic 3D representations.

This repository is a fork of [g3n/engine](https://github.com/g3n/engine).

## Features

* Cross-platform: Windows, Linux, and macOS
* Integrated GUI with many widgets
* Hierarchical scene graph
* 3D spatial audio via [OpenAL](https://www.openal.org/)
* Real-time lighting with ambient, directional, point, and spot lights
* Physically-based rendering shaders
* Geometry generators for box, sphere, cylinder, torus, etc...
* Geometries support morph targets and multimaterials
* Support for animated sprites based on sprite sheets
* Perspective and orthographic cameras
* Text image generation and support for TrueType fonts
* Image textures can be loaded from GIF, PNG or JPEG files
* Animation framework for position, rotation, and scale of objects
* Support for user-created GLSL shaders
* Support for HiDPI displays

## Dependencies

**Go 1.23+** is required.

An **OpenGL driver** and a **GCC-compatible C compiler** is required.

See below for OS specific requirements.

### Ubuntu/Debian-like

```shell
sudo apt-get install xorg-dev libgl1-mesa-dev libopenal1 libopenal-dev libvorbis0a libvorbis-dev libvorbisfile3
```

### Fedora

```shell
sudo dnf -y install xorg-x11-proto-devel mesa-libGL mesa-libGL-devel openal-soft openal-soft-devel libvorbis libvorbis-devel glfw-devel libXi-devel libXxf86vm-devel
```

### CentOS 7

```shell
# enable the EPEL repository
sudo yum -y install https://dl.fedoraproject.org/pub/epel/epel-release-latest-7.noarch.rpm
sudo yum -y install xorg-x11-proto-devel mesa-libGL mesa-libGL-devel openal-soft openal-soft-devel libvorbis libvorbis-devel glfw-devel libXi-devel libXxf86vm-devel
```

### Arch

```shell
sudo pacman -S base-devel xorg-server mesa openal libvorbis
```

### Void

```shell
sudo xbps-install git xorg-server-devel base-devel libvorbis-devel libvorbis libXxf86vm-devel libXcursor-devel libXrandr-devel libXinerama-devel libopenal libopenal-devel libglvnd-devel
```

### Windows

Use the [mingw-w64](https://mingw-w64.org) toolchain by
downloading [this file](https://sourceforge.net/projects/mingw-w64/files/Toolchains%20targetting%20Win64/Personal%20Builds/mingw-builds/8.1.0/threads-posix/seh/x86_64-8.1.0-release-posix-seh-rt_v6-rev0.7z).

Download the necessary [audio DLLs](audio/windows/bin) and add them to your PATH, or build the DLLs yourself with
instructions [here](audio/windows).

### macOS

```shell
brew install libvorbis openal-soft
```
