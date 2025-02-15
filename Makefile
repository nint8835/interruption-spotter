css:
	tailwindcss -i pkg/server/static/input.css -o pkg/server/static/style.css 

prod-css:
	tailwindcss -i pkg/server/static/input.css -o pkg/server/static/style.css --minify
