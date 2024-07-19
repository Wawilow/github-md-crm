# https://gitpython.readthedocs.io/en/stable/tutorial.html
# pip3 install gitpython
# pip3 install pygithub3

from github import Github
from pygithub3 import Github
import urllib
import base64

token = ""


# read file data
def read_file(repo, path, git_token, branch="main"):
    g = Github(git_token)
    repo = g.get_repo(repo)
    content_encoded = repo.get_contents(urllib.parse.quote(path), ref=branch).content
    content = base64.b64decode(content_encoded)
    return content


def update_file(repo, path, git_token, branch="main"):
    g = Github(git_token)
    repo = g.get_repo(repo)
    found=False
    try:
        found = repo.get_contents(
            path=urllib.parse.quote(path),
        )
    except:
        pass

    if found == False:
        return repo.create_file(
            path=urllib.parse.quote(path),
            message="test commit",
            content="-\n--\n",
            branch=branch,
        )
    else:
        return repo.update_file(
            path=urllib.parse.quote(path),
            message="test commit",
            sha=found.sha,
            content="-\n--2\n",
            branch=branch,
        )


print(update_file(repo="", path="", git_token=token, branch="main"))
