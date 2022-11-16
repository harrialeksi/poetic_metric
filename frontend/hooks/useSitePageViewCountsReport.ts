import { useMemo } from "react";
import useSWR, { SWRResponse } from "swr";
import { hydrateSitePageViewCountsReport } from "../helpers";
import { useReportQueryParams } from "./useReportQueryParams";

type HydratedSwrResponse = SWRResponse<SitePageViewCountsReport> & {
  hydratedData?: HydratedSitePageViewCountsReport,
};

export function useSitePageViewCountsReport(): HydratedSwrResponse {
  const reportQueryParams = useReportQueryParams();
  const swrResponse = useSWR<SitePageViewCountsReport>(`/site-reports/page-view-counts?${reportQueryParams}`);

  const hydratedData = useMemo<HydratedSitePageViewCountsReport | undefined>(() => {
    if (swrResponse.data === undefined) {
      return undefined;
    }

    return hydrateSitePageViewCountsReport(swrResponse.data);
  }, [swrResponse.data]);

  return { ...swrResponse, hydratedData };
}
