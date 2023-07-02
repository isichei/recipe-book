from fastapi import FastAPI, Request
from fastapi.responses import HTMLResponse, FileResponse
from fastapi.templating import Jinja2Templates
from fastapi.staticfiles import StaticFiles

app = FastAPI()

app.mount("/static", StaticFiles(directory="static"), name="static")
app.mount("/thumbnails", StaticFiles(directory="thumbnails"), name="thumbnails")

templates = Jinja2Templates(directory="templates")

RECIPES_DB = [
    {
        "uid": "chicken-dhansak-recipe",
        "title": "Chicken Dhansak",
        "description": "A chicken dhansak recipe from BBC good foods",
    },
    {
        "uid": "christmas-roast-potatoes",
        "title": "Jamie Oliver Roast Potatoes",
        "description": "A jamie oliver roast potato recipe usually used at Christmas",
    },
]

def _search_recipes(text:str) -> list[dict]:
    if not text:
        return RECIPES_DB[:10]

    for item in RECIPES_DB:
        matched = [
            item
            for item in RECIPES_DB
            if any(
                [word in item["description"].lower() for word in text.lower().split()]
            )
        ]
        return matched
    

@app.get("/", response_class=HTMLResponse)
async def home(request: Request):
    context = {
        "request": request,
        "results": _search_recipes(""),
    }
    return templates.TemplateResponse("home.html", context=context)


@app.get("/search-recipes", response_class=HTMLResponse)
async def search_recipes(request: Request, text: str):
    context = {
        "request": request,
        "results": _search_recipes(text),
    }
    return templates.TemplateResponse("search_results.html", context=context)


# @app.get("/thumbnail/{name}", response_class=FileResponse)
# async def thumbnail(name:str):
#     return FileResponse(f"thumbnails/{name}")
