import { test, expect } from '../fixtures';

test.describe('Tareas', () => {
  test('editar la fecha límite', async ({ page, tasks }) => {
    await tasks([{ title: 'Entregar el informe' }, { title: 'Revisar el presupuesto' }]);
    await page.reload();

    // 1. Sin fecha límite, la lista no la muestra
    await expect(page.getByText('Fecha límite:')).toHaveCount(0);

    // 2. Pone una fecha límite a «Entregar el informe»
    await page.getByRole('link', { name: 'Editar tarea: Entregar el informe' }).click();
    await page.getByLabel('Fecha límite').fill('2026-10-25');
    await page.getByRole('button', { name: 'Guardar cambios' }).click();

    // 3. La lista la muestra en su tarea, también al recargar
    await expect(page.getByRole('status')).toHaveText('Tarea actualizada.');
    const task = page.getByRole('listitem').filter({ hasText: 'Entregar el informe' });
    await expect(task.getByText('Fecha límite: 25/10/2026')).toBeVisible();
    await page.reload();
    await expect(task.getByText('Fecha límite: 25/10/2026')).toBeVisible();
    await expect(page.getByText('Fecha límite:')).toHaveCount(1);

    // 4. La quita dejando el campo vacío
    await page.getByRole('link', { name: 'Editar tarea: Entregar el informe' }).click();
    await expect(page.getByLabel('Fecha límite')).toHaveValue('2026-10-25');
    await page.getByLabel('Fecha límite').fill('');
    await page.getByRole('button', { name: 'Guardar cambios' }).click();
    await expect(page.getByRole('status')).toHaveText('Tarea actualizada.');
    await expect(page.getByText('Fecha límite:')).toHaveCount(0);
  });
});
