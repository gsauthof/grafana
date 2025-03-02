import React from 'react';
import { render, screen } from '@testing-library/react';
import { ColorPickerPopover } from './ColorPickerPopover';
import { createTheme } from '@grafana/data';
import userEvent from '@testing-library/user-event';

describe('ColorPickerPopover', () => {
  const theme = createTheme();

  it('should be tabbable', () => {
    render(<ColorPickerPopover color={'red'} onChange={() => {}} />);
    const color = screen.getByRole('button', { name: 'red color' });
    const customTab = screen.getByRole('button', { name: 'Custom' });

    userEvent.tab();
    expect(customTab).toHaveFocus();

    userEvent.tab();
    expect(color).toHaveFocus();
  });

  describe('rendering', () => {
    it('should render provided color as selected if color provided by name', () => {
      render(<ColorPickerPopover color={'green'} onChange={() => {}} />);
      const color = screen.getByRole('button', { name: 'green color' });
      const colorSwatchWrapper = screen.getAllByTestId('data-testid-colorswatch');

      expect(color).toBeInTheDocument();
      expect(colorSwatchWrapper[0]).toBeInTheDocument();

      userEvent.click(colorSwatchWrapper[0]);
      expect(color).toHaveStyle('box-shadow: inset 0 0 0 2px #73BF69, inset 0 0 0 4px #000000');
    });
  });

  describe('named colors support', () => {
    const onChangeSpy = jest.fn();

    it('should pass hex color value to onChange prop by default', () => {
      render(<ColorPickerPopover color={'red'} onChange={onChangeSpy} />);
      const color = screen.getByRole('button', { name: 'red color' });
      userEvent.click(color);

      expect(onChangeSpy).toBeCalledTimes(1);
      expect(onChangeSpy).toBeCalledWith(theme.visualization.getColorByName('red'));
    });

    it('should pass color name to onChange prop when named colors enabled', () => {
      render(<ColorPickerPopover color={'red'} enableNamedColors onChange={onChangeSpy} />);
      const color = screen.getByRole('button', { name: 'red color' });
      userEvent.click(color);

      expect(onChangeSpy).toBeCalledTimes(2);
      expect(onChangeSpy).toBeCalledWith(theme.visualization.getColorByName('red'));
    });
  });
});
