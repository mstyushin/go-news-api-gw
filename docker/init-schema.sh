#!/usr/bin/env bash

psql -v ON_ERROR_STOP=1 --username gonews --dbname news <<-SQL
DROP TABLE IF EXISTS comments;
CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    parent_id INTEGER DEFAULT 0,
    FOREIGN KEY (parent_id) REFERENCES comments(id),
    article_id INTEGER DEFAULT 0,
    author TEXT  NOT NULL,
    text TEXT NOT NULL,
    pub_time INTEGER DEFAULT 0
);
CREATE INDEX idx_comments_pub_time ON comments (pub_time);
DROP TABLE IF EXISTS articles, rss_feeds;
CREATE TABLE rss_feeds (
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL,
    UNIQUE(url)
);
CREATE TABLE articles (
    id SERIAL PRIMARY KEY,
    rss_feed_id INTEGER REFERENCES rss_feeds(id) NOT NULL,
    title TEXT  NOT NULL,
    content TEXT NOT NULL,
    link TEXT NOT NULL,
    pub_time INTEGER DEFAULT 0,
    UNIQUE(link)
);
CREATE INDEX idx_articles_pub_time ON articles (pub_time);
SQL
