import 'package:flutter/material.dart';

class CardActions extends StatelessWidget {
  final VoidCallback? onView;
  final VoidCallback? onEdit;
  final VoidCallback? onDelete;
  final VoidCallback? onMore;
  final VoidCallback? onLive;
  final bool isEditDisabled;
  final bool isDeleteDisabled;
  final bool showDeleteLoader;
  const CardActions({
    super.key,
    this.onView,
    this.onEdit,
    this.onDelete,
    this.onMore,
    this.onLive,
    this.isEditDisabled = false,
    this.isDeleteDisabled = false,
    this.showDeleteLoader = false,
  });

  ButtonStyle _primaryStyle(BuildContext context) {
    return FilledButton.styleFrom(
      minimumSize: const Size(0, 36),
      padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 10),
      tapTargetSize: MaterialTapTargetSize.shrinkWrap,
      visualDensity: const VisualDensity(horizontal: -2, vertical: -2),
    );
  }

  ButtonStyle _tonalStyle(BuildContext context) {
    return FilledButton.styleFrom(
      minimumSize: const Size(0, 36),
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
      tapTargetSize: MaterialTapTargetSize.shrinkWrap,
      visualDensity: const VisualDensity(horizontal: -2, vertical: -2),
    );
  }

  ButtonStyle _dangerStyle(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    return TextButton.styleFrom(
      minimumSize: const Size(0, 36),
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 8),
      tapTargetSize: MaterialTapTargetSize.shrinkWrap,
      foregroundColor: cs.error,
      visualDensity: const VisualDensity(horizontal: -2, vertical: -2),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Wrap(
      spacing: 8,
      runSpacing: 8,
      children: [
        if (onView != null)
          FilledButton.icon(
            style: _primaryStyle(context),
            onPressed: onView,
            icon: const Icon(Icons.visibility_outlined, size: 18),
            label: const Text('Ver'),
          ),
        if (onEdit != null)
          FilledButton.tonalIcon(
            style: _tonalStyle(context),
            onPressed: isEditDisabled ? null : onEdit,
            icon: const Icon(Icons.edit_outlined, size: 18),
            label: const Text('Editar'),
          ),
        if (onDelete != null)
          TextButton.icon(
            style: _dangerStyle(context),
            onPressed: (isDeleteDisabled || showDeleteLoader) ? null : onDelete,
            icon: showDeleteLoader
                ? SizedBox(
                    width: 18,
                    height: 18,
                    child: CircularProgressIndicator(
                      strokeWidth: 2.2,
                      color: Theme.of(context).colorScheme.error,
                    ),
                  )
                : const Icon(Icons.delete_outline, size: 18),
            label: Text(showDeleteLoader ? 'Eliminando...' : 'Eliminar'),
          ),
        if (onMore != null)
          IconButton(
            onPressed: onMore,
            icon: const Icon(Icons.more_vert),
          ),
        if (onLive != null)
          FilledButton.icon(
            style: _primaryStyle(context),
            onPressed: onLive,
            icon: const Icon(Icons.live_tv, size: 18),
            label: const Text('Votación en vivo'),
          ),
      ],
    );
  }
}
