import { useRouter } from "next/router";
import React, { useContext, useEffect, useMemo } from "react";
import { Container, Spinner } from "react-bootstrap";
import { CanonicalLink, Layout, SiteReports, Title } from "../components";
import { ToastsContext } from "../contexts";
import { usePublicSite } from "../hooks";

export default function PublicSiteReports() {
  const router = useRouter();
  const { addToast } = useContext(ToastsContext);

  const domain = useMemo<null | string | undefined>(() => {
    if (!router.isReady) {
      return undefined;
    }

    return router.query.d?.toString() || null;
  }, [router.isReady, router.query.d]);

  const { data: site, error: siteError } = usePublicSite(domain || undefined);

  useEffect(() => {
    if (domain === null) {
      router.replace("/");
    }
  }, [domain, router]);

  useEffect(() => {
    if (siteError !== undefined) {
      addToast({ body: siteError.message || "An error has occurred.", variant: "danger" });

      router.replace("/");
    }
  }, [addToast, router, siteError]);

  return (
    <Layout kind="website">
      <Container className="d-flex flex-column flex-grow-1 py-5">
        {site === undefined ? (
          <Spinner className="m-auto" variant="primary" />
        ) : (
          <>
            {!!domain ? (
              <CanonicalLink path={`/s?d=${domain}`} />
            ) : null}

            <Title>{`Reports for ${site.name}`}</Title>

            <SiteReports site={site} />
          </>
        )}
      </Container>
    </Layout>
  );
}
