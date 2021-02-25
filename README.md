<p align="center">
    <img src="https://omniskop.de/firefly/images/logo.png" width="20%" alt="Firefly logo. A vinyl disk with red, green and blue light reflecting off of it.">
</p>

# Firefly
Firefly is a program for animating LED-Strips to music.

The editor is based on the idea of a graphics editor on which the y axis represents time and the x axis is mapped to the LED-Strip. On this canvas you can place and manipulate different shapes to create complex animations.

<p align="center">
    <img src="https://omniskop.de/firefly/images/editor.png" width="90%" alt="A screenshot of the Firefly editor. Showing multiple stripes that are animated to move over the led strip over time.">
</p>

At the top you can see the "needle" which shows the image of the LED-Strip at the current point in time. The editor can stream these animations live via UDP to an LED-Strip connected to a WLED-based micro controller. (It is easy to add more protocols for more compatibility.) The number of LEDs can be configured and even complex LED-Strip arrangements with varying pixel density are supported. Due to the nature of this interface concept only 1D animations are supported and that is unlikely to change.  
Firefly will play the animation and music in sync and for more precision during editing playback at slower speeds is supported.

Firefly is still very much work in progress. It is currently running on Windows, Linux and macOS. 