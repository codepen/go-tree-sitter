#include "tree_sitter/parser.h"
#include "tree_sitter/alloc.h"

#define SHOW_DEBUG 0
#define DEBUG(...) if (SHOW_DEBUG) { fprintf(stderr, __VA_ARGS__); }
enum TokenType {
  COMMENT,
  STRING,
};

/* We just need to hold the last token and nothing else */
struct Scanner {
  int32_t last_token;
};

void *tree_sitter_jinja2_external_scanner_create() {
  return NULL;
}

unsigned tree_sitter_jinja2_external_scanner_serialize(void *payload, char *buffer) {
  return 0;
}

void tree_sitter_jinja2_external_scanner_deserialize(void *payload, const char *buffer, unsigned length) {
}

void tree_sitter_jinja2_external_scanner_destroy(void *payload) {
}

static inline void skip(TSLexer *lexer) {
  lexer->advance(lexer, true);
}

static inline void advance(TSLexer *lexer) {
  lexer->advance(lexer, false);
}

bool tree_sitter_jinja2_external_scanner_scan(void *payload, TSLexer *lexer, const bool *valid_symbols) {
  bool has_match = false;
  /* skip whitespace but not for comment as '{ #' is invalid and
   * would match. */
  while (!valid_symbols[COMMENT] &&
         (lexer->lookahead == ' ' ||
          lexer->lookahead == '\t' ||
          lexer->lookahead == '\n' ||
          lexer->lookahead == '\r')) {
    skip(lexer);
  }

  /* Handle strings including \ to escape the symbols inside the string */
  if (valid_symbols[STRING]) {
    DEBUG("STRING next char: %c\n", lexer->lookahead);
    bool has_start_double_quote = lexer->lookahead == '"';
    bool has_start_single_quote = lexer->lookahead == '\'';
    if (has_start_single_quote || has_start_double_quote) {
      DEBUG("has_start_single_quote: %d\n", has_start_single_quote);
      DEBUG("has_start_double_quote: %d\n", has_start_double_quote);
      advance(lexer);
      bool has_escape = false;
      bool has_end_double_quote = lexer->lookahead == '"';
      bool has_end_single_quote = lexer->lookahead == '\'';
      while (!lexer->eof(lexer)) {
        DEBUG("loop\n");
        has_end_double_quote = lexer->lookahead == '"';
        has_end_single_quote = lexer->lookahead == '\'';
        has_escape = lexer->lookahead == '\\';
        if (has_escape)
          advance(lexer);
        if ((has_start_double_quote && has_end_double_quote)
            || (has_start_single_quote && has_end_single_quote)) {
          DEBUG("has_end_double_quote: %d\n", has_end_double_quote);
          DEBUG("has_end_single_quote: %d\n", has_end_single_quote);
          advance(lexer);
          lexer->result_symbol = STRING;
          lexer->mark_end(lexer);
          return true;
        }
        DEBUG("has_escape: %d\n", has_escape);
        advance(lexer);
      }
    } else {
      lexer->result_symbol = 0;
      DEBUG("Missing\n");
      return false;
    }
  } else if (lexer->lookahead == '#' && valid_symbols[COMMENT]) {
    DEBUG("COMMENT\n");
    bool has_open_comment = true;
    bool has_hash_tag = false;
    advance(lexer);
    while (!lexer->eof(lexer)) {
      if (lexer->lookahead != '}')
        has_hash_tag = false;
      if (lexer->lookahead == '#')
        has_hash_tag = true;
      if (lexer->lookahead == '}' && has_open_comment && has_hash_tag) {
        lexer->result_symbol = COMMENT;
        lexer->mark_end(lexer);
        return true;
      }
      advance(lexer);
    }
  } else{
    lexer->result_symbol = 0;
  }
  return false;
}



