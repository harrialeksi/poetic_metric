import classNames from "classnames";
import Link from "next/link";
import React, { useId } from "react";
import { Dropdown } from "react-bootstrap";
import styles from "./Menu.module.scss";
import { Search } from "./Search";

export type MenuProps = Overwrite<Omit<React.PropsWithoutRef<JSX.IntrinsicElements["div"]>, "children">, {
  article: DocsArticle;
  categories: Array<DocsCategory>;
}>;

export function Menu({ article: articleFromProps, categories, className, ...props }: MenuProps) {
  const dropdownToggleId = useId();

  return (
    <>
      <div className={classNames("d-md-none p-3 sticky-top", styles.smallMenu, className)} {...props}>
        <Dropdown>
          <Dropdown.Toggle className="d-block m-auto" id={dropdownToggleId}>
            <i className="bi bi-journal-text me-2" />

            Show docs menu
          </Dropdown.Toggle>

          <Dropdown.Menu>
            <div className="p-2">
              <Search />
            </div>

            <div className={classNames("gap-2 pt-1 vstack", styles.links)}>
              {categories.map((category) => (
                <div key={category.slug}>
                  <h6 className="px-3">{category.title}</h6>

                  {category.articles.map((article) => (
                    <Dropdown.Item
                      as={Link}
                      href={`/docs/${category.slug}/${article.slug}`}
                      key={article.slug}
                    >
                      {article.title}
                    </Dropdown.Item>
                  ))}
                </div>
              ))}
            </div>
          </Dropdown.Menu>
        </Dropdown>
      </div>

      <div
        className={classNames(
          "d-none d-md-block flex-grow-0 flex-shrink-0 overflow-auto px-3 py-4 sticky-top w-16rem",
          styles.menu,
          className,
        )}
        {...props}
      >
        <nav className="gap-3 vstack">
          <Search />

          {categories.map((category) => (
            <div key={category.slug}>
              <h6>{category.title}</h6>

              <div className="gap-1 vstack">
                {category.articles.map((article) => (
                  <Link
                    className={classNames(
                      "d-block text-decoration-none",
                      articleFromProps.category.slug === category.slug && articleFromProps.slug === article.slug
                        ? "text-primary"
                        : "text-decoration-underline-hover text-dark",
                    )}
                    href={`/docs/${category.slug}/${article.slug}`}
                    key={article.slug}
                  >
                    {article.title}
                  </Link>
                ))}
              </div>
            </div>
          ))}
        </nav>
      </div>
    </>
  );
}
