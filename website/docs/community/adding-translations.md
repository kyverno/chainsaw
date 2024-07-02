# Translations

Chainsaw, developed by Kyverno, benefits greatly from international community contributions. To make Chainsaw even more globally accessible, we encourage you to help us by adding missing translations or new languages.

## Before Creating an Issue

Translations change frequently, so we want to ensure you don't duplicate work. Before adding translations, please check the following:

### Check Language Availability

Chainsaw may already support your language. Check the list of [supported languages] to see if your language is available or needs improvements:

- __Your language is already supported__ – Check for missing translations and click the link under your language to add them.
- __Your language is missing__ – Help us add support for your language! Read on to learn how.

<!-- [supported languages]: ../setup/changing-the-language.md#site-language -->

### Search Our Issue Tracker

Another user might have already created an issue for missing translations in your language. To avoid duplicating work, search the [issue tracker] beforehand.

[issue tracker]: https://github.com/kyverno/chainsaw/issues

---

If Chainsaw doesn't already support your language, you can add new translations by following the issue template below.

## Issue Template

We have created an issue template to make contributing translations simple. It consists of the following parts:

- [Title]
- [Translations]
- [Country Flag] <small>optional</small>
- [Checklist]

[Title]: #title
[Translations]: #translations
[Country Flag]: #country-flag
[Checklist]: #checklist

### Title

When updating an existing language, leave the title as is. For a new language, replace `...` in the pre-filled title with the name of your language.

| <!-- --> | Example  |
| -------- | -------- |
| :material-check:{ style="color: #4DB6AC" } __Clear__ | Add translations for German
| :material-close:{ style="color: #EF5350" } __Unclear__ | Add translations ...
| :material-close:{ style="color: #EF5350" } __Useless__ | Help

### Translations

If a translation contains an :arrow_left: icon on the right side, it is missing. Translate this line and remove the :arrow_left: icon. If unsure about specific lines, leave them for others to complete. Double-check the context by looking at our [English translations].

<!-- [English translations]: https://github.com/kyverno/chainsaw/tree/master/src/partials/languages/en.html -->

### Country Flag <small>optional</small> {#country-flag}

For better visibility, our list of [supported languages] includes country flags next to the language names. Help us select a flag for your language by adding the shortcode for the country flag. Use the [emoji search] to find all available shortcodes.

!!! question "What if my flag is not available?"

    [Twemoji] provides flag emojis for 260 countries. Subdivisions of countries, such as states or regions, are not supported. If adding translations for a subdivision, choose the most appropriate available flag.

[Twemoji]: https://twemoji.twitter.com/
<!-- [emoji search]: ../reference/icons-emojis.md#search -->

> __Why this might be helpful__: Adding a country flag helps others find the language in the list of supported languages faster. If your country's flag is not supported by [Twemoji], help us choose an alternative.

### Checklist

Thanks for following the guide and helping us add new translations to Chainsaw – you are almost done. The checklist ensures you have read this guide and provided all necessary information for us to integrate your contribution.

__We'll take it from here.__

---

## Attribution

If you submit a translation using the template above, you will be __credited as a co-author__ in the commit, so you don't need to open a pull request. Your significant contribution makes Chainsaw more accessible to people worldwide. Thank you!
