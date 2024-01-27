import { PRODUCT_STATS } from '@/lib/mock/product';
import { render, waitFor } from '@testing-library/vue';
import { Mock } from 'vitest';
import MainPage from '@/pages/Main/index.vue'

global.fetch = vi.fn(() =>
  Promise.resolve({
    json: () => 
      Promise.resolve(PRODUCT_STATS)
  })
) as Mock

const renderComponent = () =>  render(MainPage)

describe('Main page', () => {
  afterEach(() => {
    vi.restoreAllMocks()
  })
  test('Should render page correctly',async () => {
    const { container, findAllByTestId } = renderComponent()
    await waitFor(() => {
      expect(global.fetch)
        .toHaveBeenLastCalledWith("http://localhost:8080/api/product-sale-stats")
    })

    const productRows = await findAllByTestId('product-sales-rows')
    expect(productRows).toBeDefined()
    expect(container).toMatchSnapshot()
  })
})
