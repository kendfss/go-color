color
===

Package color implements some simple RGB/HSL color conversions (+-0.02% of error) for golang.

Adapted from
http://code.google.com/p/closure-library/source/browse/trunk/closure/goog/color/color.js
and algorithms on easyrgb.com.

To maintain accuracy between conversions we use floats in the color types.
If you are storing lots of colors and care about memory use you might want
to use something based on byte types instead.

Also, color types don't verify their validity before converting. If you do
something like RGB{10,20,30}.ToHSL() the results will be undefined. You can 
use the `New` initializer for values between 0 and 255. Otherwise, All 
values must be between 0 and 1.

By Brandon Thomson
