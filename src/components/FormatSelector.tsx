import { InlineFormLabel, Select } from '@grafana/ui';
import React from 'react';
import { isDataQuery } from './../app/utils';
import { INFINITY_RESULT_FORMATS } from './../constants';
import { Components } from './../selectors';
import type { InfinityQuery, InfinityQueryFormat } from './../types';
interface FormatSelectorProps {
  query: InfinityQuery;
  onChange: (e: InfinityQuery) => void;
  onRunQuery: () => void;
}
export const FormatSelector = (props: FormatSelectorProps) => {
  const { query, onChange, onRunQuery } = props;
  if (!(isDataQuery(query) || query.type === 'uql')) {
    return <></>;
  }
  const onFormatChange = (format: InfinityQueryFormat) => {
    onChange({ ...query, format });
    onRunQuery();
  };
  const getFormats = () => {
    if (query.type === 'json') {
      return INFINITY_RESULT_FORMATS;
    } else if (query.type === 'uql') {
      return INFINITY_RESULT_FORMATS.filter((f) => f.value === 'table' || f.value === 'timeseries' || f.value === 'dataframe');
    } else {
      return INFINITY_RESULT_FORMATS.filter((f) => f.value !== 'as-is');
    }
  };
  return (
    <>
      <InlineFormLabel width={4} className={`query-keyword`}>
        {Components.QueryEditor.Format.Label.Text}
      </InlineFormLabel>
      <div title={Components.QueryEditor.Format.Dropdown.PlaceHolder.Title} style={{ marginRight: '5px' }} data-testid="infinity-query-format-selector">
        <Select className="min-width-12 width-12" value={query.format} options={getFormats()} onChange={(e) => onFormatChange(e.value as InfinityQueryFormat)} menuShouldPortal={true} />
      </div>
    </>
  );
};
