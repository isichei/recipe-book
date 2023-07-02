from fastapi import FastAPI, Request
from fastapi.responses import HTMLResponse
from fastapi.templating import Jinja2Templates

app = FastAPI()
templates = Jinja2Templates(directory="templates")

RECIPES_DB = [
    {
        "name": "chicken-dhansak-recipe",
        "description": "A chicken dhansak recipe from BBC good foods",
        "source": "recipes/chicken-dhansak-recipe.pdf",
    },
    {
        "name": "christmas-roast-potatoes",
        "description": "A jamie oliver roast potato recipe usually used at xmas",
        "source": "recipes/jamie-oliver-roast-potatoes.pdf",
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
def home(request: Request):
    context = {
        "request": request,
        "results": _search_recipes(""),
    }
    return templates.TemplateResponse("home.html", context=context)


@app.get("/search-recipes", response_class=HTMLResponse)
def search_recipes(request: Request, text: str):
    context = {
        "request": request,
        "results": _search_recipes(text),
    }
    return templates.TemplateResponse("search_results.html", context=context)
