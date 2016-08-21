/*
  Package `servefiles` provides a static asset handler for serving files such as images, stylesheets and
  javascript code. Care is taken to set headers such that the assets will be efficiently cached by browsers and proxies.

      assets := servefiles.AssetHandler(0, "./assets/", time.Hour)

  Assets is an http.Handler and can be used alongside your other handlers.

  Gzipping

  During the preparation of your web assets, all text files (CSS, JS etc) should be accompanied by their gzipped
  equivalent. The Assets handler will first check for the gzipped file, which it will serve if present. Otherwise
  it will serve the 'normal' file.

  This has many benefits: fewer bytes are read from the disk, a smaller memory footprint is needed in the server,
  fewer bytes are sent across the network, etc.

  You should not attempt to gzip already-compressed files, such as JPEG, SVGZ, etc.

  Very small files gain little from compression because they may be small enough to fit within a single TCP packet,
  so don't bother with them. (They might even grow in size.)

  Cache Control

  The 'far-future' technique can and should be used. Set a long expiry time, e.g. time.Hour * 24 * 3650

  No in-memory caching is performed server-side. This is less necessary due to far-future caching being
  supported, but might be added in future.

  Example Usage

  Choose any expiry age you like, although ten years seems to work well.
  To serve files with a ten-year expiry, this creates a suitably-configured handler:

      assets := servefiles.AssetHandler(1, "./assets/", 10 * 365 * 24 * time.Hour)

  Notice the first parameter is 1 instead of the default, 0. So the first segment of the URL path is discarded.
  This means, for example, that the URL

      http://example.com/e3b1cf/css/style1.css

  maps to the asset files

      ./assets/css/style1.css
      ./assets/css/style1.css.gz

  without the e3b1cf segment. The benefit of this is that you can use a unique number in that segment (chosen for
  example each time your server starts). Each time that number changes, browsers will see the asset files as
  being new, and later will drop old versions from their cache regardless of their ten-year lifespan.

  So you get the ten-year lifespan combined with being able to push out changed assets when you need to.
*/
package servefiles

