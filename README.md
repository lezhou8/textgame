## textgame

Transform a game description into a interactive text adventure game.

```
Welcome to textgame. What would you like to do?

1. Download textgame and start playing interactive text adventure games with it.
2. Do anything else.
(Choose 1...2)
```

### Creating game descriptions

1. Create a .yaml file.
2. Create a root-level key called `entry_dialog`. All dialogs require a child key `text`. Its value is a string that will display when the player reaches the dialog. There are three kinds of dialogs:
    - Plain text. This dialog has a `plain_text` child key. `plain_text` has a child key `next_dialog` whose value specifies the next dialog/root-level key that is displayed.
    - Multi-choice. This dialog has a `multi_choice` child key. The value of `multi_choice` is a sequence of mappings with keys `text` and `next_dialog`, each with values displaying the option and next dialog if that option is selected respectively.
    - Text-box. This dialog has a `text_box` child key. The value of `text_box` is a sequence of mappings with keys `text` and `next_dialog`, each with values for a possible input and the next dialog if that input is entered respectively.
3. Create a root-level key called `last_dialog`. This dialog only needs a `text` key value pair.

Examples of game descriptions are provided in the `examples` directory.

### Build

```
go build
```

### Play

```
./textgame <game_description>
```
