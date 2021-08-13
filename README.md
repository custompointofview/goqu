# GoQu

A Golang CLI tool in search of quotes Quotes

![Alt Text](./assets/profile/goku.gif)

## Goal

I wanted a simple CLI tool to get quotes based on preferences.

## Current Sources

> Additional source will be added in the near future

-   QuoteGarden: [GitHub Repo](https://github.com/pprathameshmore/QuoteGarden)

## Packages used

Awesome packages used for CLI interactions and rendering:

-   [PTerm](http://github.com/pterm/pterm)
-   [PromptUI](github.com/manifoldco/promptui)

## Notes on implementation...

I've tried something new this time regarding error handling in Go, a more Pythonic approach as I treated errors like throwing exceptions.

> Conclusion: the fact that this isn't idiomatic Go makes it a rather difficult to manage pattern. It doesn't seem like a good idea for the long term on a bigger project. The `return` statement is usually the issue, even though it's a far less hassle than handling the error every time it is returned from a function.
