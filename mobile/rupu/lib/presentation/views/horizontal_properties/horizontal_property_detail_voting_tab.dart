part of 'horizontal_property_detail_view.dart';

class _VotingTab extends GetWidget<HorizontalPropertyVotingController> {
  final String controllerTag;
  const _VotingTab({required this.controllerTag});

  @override
  String? get tag => controllerTag;

  @override
  Widget build(BuildContext context) {
    return Obx(() {
      final isLoading = controller.isLoading.value;
      final error = controller.errorMessage.value;
      final groups = List<HorizontalPropertyVotingGroup>.of(controller.groups);
      final bottomPadding = MediaQuery.viewInsetsOf(context).bottom;
      return RefreshIndicator(
        onRefresh: controller.refresh,
        child: Stack(
          children: [
            ListView(
              physics: const AlwaysScrollableScrollPhysics(),
              padding: const EdgeInsets.fromLTRB(16, 16, 16, 24),
              children: [
                SummaryHeader(
                  title: 'Grupos de votación',
                  subtitle: groups.isEmpty
                      ? 'Sin grupos registrados'
                      : '${groups.length} grupos disponibles',
                  showProgress: isLoading,
                  onRefresh: controller.refresh,
                ),
                const SizedBox(height: 16),
                if (error != null) ...[
                  _InlineError(message: error),
                  const SizedBox(height: 16),
                ],
                if (!isLoading && groups.isEmpty)
                  const _EmptyState(
                    icon: Icons.how_to_vote_outlined,
                    title: 'No hay votaciones registradas.',
                    subtitle:
                        'Cuando se creen nuevos procesos de votación aparecerán aquí.',
                  )
                else ...[
                  ...groups
                      .map(
                        (group) => Padding(
                          padding: const EdgeInsets.only(bottom: 16),
                          child: _VotingGroupCard(
                            controllerTag: controllerTag,
                            group: group,
                            onOpenAttendance: () {
                              final state = GoRouterState.of(context);
                              final segments = state.uri.pathSegments;
                              final page =
                                  segments.length > 1 ? segments[1] : '0';
                              final propertyId = controller.propertyId;
                              final path =
                                  '/home/$page/horizontal-properties/$propertyId/voting/${group.id}/attendance';
                              context.push(path, extra: group);
                            },
                          ),
                        ),
                      )
                      .toList(),
                ],
              ],
            ),
            Positioned(
              right: 24,
              bottom: 24 + bottomPadding,
              child: _AddVotingGroupFab(controllerTag: controllerTag),
            ),
          ],
        ),
      );
    });
  }
}

class _VotingGroupCard extends StatelessWidget {
  final HorizontalPropertyVotingGroup group;
  final VoidCallback onOpenAttendance;
  final String controllerTag;
  const _VotingGroupCard(
      {required this.group,
      required this.onOpenAttendance,
      required this.controllerTag});

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    final tt = Theme.of(context).textTheme;

    final (bgChip, fgChip, labelChip) = group.isActive
        ? (cs.secondaryContainer, cs.onSecondaryContainer, 'ACTIVO')
        : (cs.errorContainer, cs.onErrorContainer, 'INACTIVO');

    return Material(
      color: Colors.transparent,
      child: InkWell(
        borderRadius: BorderRadius.circular(18),
        onTap: () {
          final state = GoRouterState.of(context);
          final segments = state.uri.pathSegments;
          final page = segments.length > 1 ? segments[1] : '0';
          final propertyId = Get.find<HorizontalPropertyVotingController>(
                  tag: controllerTag)
              .propertyId;
          final path =
              '/home/$page/horizontal-properties/$propertyId/voting/${group.id}';
          context.push(path);
        },
        child: Ink(
          decoration: BoxDecoration(
            color: cs.surface,
            borderRadius: BorderRadius.circular(18),
            border: Border.all(color: cs.outlineVariant),
            boxShadow: [
              BoxShadow(
                color: Colors.black.withValues(alpha: .05),
                blurRadius: 18,
                offset: const Offset(0, 10),
              ),
            ],
          ),
          child: Padding(
            padding: const EdgeInsets.fromLTRB(18, 18, 18, 16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Container(
                      width: 44,
                      height: 44,
                      decoration: BoxDecoration(
                        color: cs.primary.withValues(alpha: .12),
                        shape: BoxShape.circle,
                      ),
                      alignment: Alignment.center,
                      child: Icon(
                        Icons.how_to_vote_outlined,
                        color: cs.primary,
                      ),
                    ),
                    const Spacer(),
                    _StatusChip(
                      label: labelChip,
                      background: bgChip,
                      foreground: fgChip,
                    ),
                  ],
                ),
                const SizedBox(height: 16),
                Text(
                  group.name,
                  style: tt.titleMedium?.copyWith(fontWeight: FontWeight.w800),
                ),
                if ((group.description?.isNotEmpty ?? false)) ...[
                  const SizedBox(height: 6),
                  Text(
                    group.description!,
                    style: tt.bodyMedium?.copyWith(color: cs.onSurfaceVariant),
                  ),
                ],
                const SizedBox(height: 16),
                Wrap(
                  spacing: 12,
                  runSpacing: 12,
                  children: [
                    _DetailLine(
                      icon: Icons.calendar_month_outlined,
                      label: 'Inicio',
                      value: _formatDate(group.votingStartDate),
                    ),
                    _DetailLine(
                      icon: Icons.event_outlined,
                      label: 'Fin',
                      value: _formatDate(group.votingEndDate),
                    ),
                    _DetailLine(
                      icon: Icons.gavel_outlined,
                      label: 'Requiere quórum',
                      value: group.requiresQuorum ? 'Sí' : 'No',
                    ),
                    _DetailLine(
                      icon: Icons.percent_outlined,
                      label: 'Quórum',
                      value: group.quorumPercentage != null
                          ? '${group.quorumPercentage}%'
                          : '--',
                    ),
                  ],
                ),
                const SizedBox(height: 16),
                _CardActions(
                  onView: onOpenAttendance,
                  onEdit: () => _openEditSheet(context, group, controllerTag),
                  onDelete: () => _confirmDelete(context, group, controllerTag),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  String _formatDate(DateTime? date) {
    if (date == null) return '--';
    final day = date.day.toString().padLeft(2, '0');
    final month = date.month.toString().padLeft(2, '0');
    final year = date.year.toString();
    return '$day/$month/$year';
  }
}

Future<void> _openCreateSheet(BuildContext context, String controllerTag) async {
  final controller =
      Get.find<HorizontalPropertyVotingController>(tag: controllerTag);
  final result = await showModalBottomSheet<bool>(
    context: context,
    backgroundColor: Colors.transparent,
    useRootNavigator: true,
    isScrollControlled: true,
    builder: (_) => _VotingGroupFormBottomSheet(
      title: 'Agregar nuevo grupo de votación',
      actionLabel: 'Crear grupo',
      onSubmit: (payload) => controller.createVotingGroup(data: payload),
    ),
  );

  if (result == true) {
    _showSnack('Grupo creado', 'El grupo de votación se registró correctamente.');
  }
}

Future<void> _openEditSheet(BuildContext context,
    HorizontalPropertyVotingGroup group, String controllerTag) async {
  final controller =
      Get.find<HorizontalPropertyVotingController>(tag: controllerTag);
  final result = await showModalBottomSheet<bool>(
    context: context,
    backgroundColor: Colors.transparent,
    useRootNavigator: true,
    isScrollControlled: true,
    builder: (_) => _VotingGroupFormBottomSheet(
      title: 'Editar grupo de votación',
      actionLabel: 'Guardar cambios',
      group: group,
      onSubmit: (payload) =>
          controller.updateVotingGroup(groupId: group.id, data: payload),
    ),
  );

  if (result == true) {
    _showSnack(
        'Grupo actualizado', 'El grupo de votación se actualizó correctamente.');
  }
}

Future<void> _confirmDelete(BuildContext context,
    HorizontalPropertyVotingGroup group, String controllerTag) async {
  final cs = Theme.of(context).colorScheme;
  final confirmed = await showDialog<bool>(
    context: context,
    builder: (dialogContext) => AlertDialog(
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(18)),
      title: const Text('Eliminar grupo de votación'),
      content: Text(
        '¿Quieres eliminar el grupo ${group.name}? Esta acción no se puede deshacer.',
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.of(dialogContext).pop(false),
          child: const Text('Cancelar'),
        ),
        FilledButton(
          style: FilledButton.styleFrom(
            backgroundColor: cs.error,
            foregroundColor: cs.onError,
          ),
          onPressed: () => Navigator.of(dialogContext).pop(true),
          child: const Text('Eliminar'),
        ),
      ],
    ),
  );

  if (confirmed != true) return;

  final controller = Get.find<HorizontalPropertyVotingController>();
  final actionResult = await controller.deleteVotingGroup(group.id);

  if (actionResult.success) {
    final message = (actionResult.message?.isNotEmpty ?? false)
        ? actionResult.message!
        : 'El grupo de votación se eliminó correctamente.';
    _showSnack('Grupo eliminado', message);
  } else {
    _showSnack(
      'No se pudo eliminar',
      actionResult.message ?? 'Inténtalo nuevamente en unos instantes.',
      isError: true,
    );
  }
}

class _VotingGroupFormBottomSheet extends StatefulWidget {
  final String title;
  final String actionLabel;
  final HorizontalPropertyVotingGroup? group;
  final Future<bool> Function(
    Map<String, dynamic> data,
  ) onSubmit;

  const _VotingGroupFormBottomSheet({
    required this.title,
    required this.actionLabel,
    required this.onSubmit,
    this.group,
  });

  @override
  State<_VotingGroupFormBottomSheet> createState() =>
      _VotingGroupFormBottomSheetState();
}

class _VotingGroupFormBottomSheetState
    extends State<_VotingGroupFormBottomSheet> {
  final _formKey = GlobalKey<FormState>();
  late final TextEditingController _nameCtrl;
  late final TextEditingController _descriptionCtrl;
  late final TextEditingController _notesCtrl;
  late final TextEditingController _quorumPercentageCtrl;
  late final TextEditingController _startDateCtrl;
  late final TextEditingController _endDateCtrl;
  bool _requiresQuorum = true;
  bool _saving = false;
  String? _errorMessage;

  @override
  void initState() {
    super.initState();
    final group = widget.group;

    _nameCtrl = TextEditingController(text: group?.name ?? '');
    _descriptionCtrl = TextEditingController(text: group?.description ?? '');
    _notesCtrl = TextEditingController(text: group?.notes ?? '');
    _quorumPercentageCtrl = TextEditingController(
        text: group?.quorumPercentage?.toString() ?? '');
    _startDateCtrl = TextEditingController(
        text: group?.votingStartDate?.toIso8601String() ?? '');
    _endDateCtrl = TextEditingController(
        text: group?.votingEndDate?.toIso8601String() ?? '');
    _requiresQuorum = group?.requiresQuorum ?? true;
  }

  @override
  void dispose() {
    _nameCtrl.dispose();
    _descriptionCtrl.dispose();
    _notesCtrl.dispose();
    _quorumPercentageCtrl.dispose();
    _startDateCtrl.dispose();
    _endDateCtrl.dispose();
    super.dispose();
  }

  Map<String, dynamic> _buildPayload() {
    final payload = <String, dynamic>{
      'name': _nameCtrl.text.trim(),
      'description': _descriptionCtrl.text.trim(),
      'notes': _notesCtrl.text.trim(),
      'quorum_percentage': double.tryParse(_quorumPercentageCtrl.text.trim()),
      'voting_start_date': _startDateCtrl.text.trim(),
      'voting_end_date': _endDateCtrl.text.trim(),
      'requires_quorum': _requiresQuorum,
    };

    payload.removeWhere((key, value) => value == null);
    return payload;
  }

  Future<void> _handleSubmit() async {
    if (_saving) return;
    final formState = _formKey.currentState;
    if (formState == null || !formState.validate()) {
      return;
    }

    FocusScope.of(context).unfocus();
    setState(() {
      _saving = true;
      _errorMessage = null;
    });

    try {
      final result = await widget.onSubmit(_buildPayload());
      if (!mounted) return;
      if (!result) {
        setState(() {
          _saving = false;
          _errorMessage =
              'No se pudo guardar el grupo de votación. Inténtalo más tarde.';
        });
        return;
      }
      Navigator.of(context).pop(result);
    } catch (_) {
      if (!mounted) return;
      setState(() {
        _saving = false;
        _errorMessage =
            'Ocurrió un error al guardar el grupo de votación. Inténtalo nuevamente.';
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    final tt = Theme.of(context).textTheme;
    final viewInsets = MediaQuery.viewInsetsOf(context);

    InputDecoration decoration(String label, {String? hint}) {
      return InputDecoration(
        labelText: label,
        hintText: hint,
        filled: true,
        fillColor: cs.surfaceContainerHighest,
        border: OutlineInputBorder(borderRadius: BorderRadius.circular(14)),
        enabledBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(14),
          borderSide: BorderSide(color: cs.outlineVariant),
        ),
        focusedBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(14),
          borderSide: BorderSide(color: cs.primary.withValues(alpha: .5)),
        ),
      );
    }

    Widget field({
      required String label,
      required TextEditingController controller,
      String? hint,
      TextInputType keyboardType = TextInputType.text,
      String? Function(String?)? validator,
    }) {
      return TextFormField(
        controller: controller,
        decoration: decoration(label, hint: hint),
        keyboardType: keyboardType,
        textInputAction: TextInputAction.next,
        validator: validator,
      );
    }

    return FractionallySizedBox(
      heightFactor: 0.95,
      child: SafeArea(
        top: false,
        child: Padding(
          padding: EdgeInsets.fromLTRB(16, 0, 16, viewInsets.bottom + 16),
          child: DecoratedBox(
            decoration: BoxDecoration(
              color: cs.surface,
              borderRadius: const BorderRadius.only(
                topLeft: Radius.circular(32),
                topRight: Radius.circular(32),
              ),
              boxShadow: [
                BoxShadow(
                  color: Colors.black.withValues(alpha: .18),
                  blurRadius: 28,
                  offset: const Offset(0, 22),
                ),
              ],
            ),
            child: ClipRRect(
              borderRadius: const BorderRadius.only(
                topLeft: Radius.circular(32),
                topRight: Radius.circular(32),
              ),
              child: Column(
                children: [
                  const SizedBox(height: 12),
                  const _SheetHandle(),
                  Padding(
                    padding: const EdgeInsets.fromLTRB(24, 10, 24, 20),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          widget.title,
                          style: tt.titleMedium?.copyWith(
                            fontWeight: FontWeight.w800,
                          ),
                        ),
                        const SizedBox(height: 6),
                        Text(
                          'Completa la información del grupo para continuar.',
                          style: tt.bodyMedium?.copyWith(
                            color: cs.onSurfaceVariant,
                          ),
                        ),
                      ],
                    ),
                  ),
                  Expanded(
                    child: SingleChildScrollView(
                      padding: const EdgeInsets.fromLTRB(24, 0, 24, 24),
                      child: Form(
                        key: _formKey,
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            ResponsiveFormGrid(
                              children: [
                                field(
                                  label: 'Nombre',
                                  controller: _nameCtrl,
                                  validator: (value) {
                                    if (value == null || value.trim().isEmpty) {
                                      return 'Ingresa el nombre del grupo';
                                    }
                                    return null;
                                  },
                                ),
                                field(
                                  label: 'Descripción',
                                  controller: _descriptionCtrl,
                                ),
                                field(
                                  label: 'Notas',
                                  controller: _notesCtrl,
                                ),
                                field(
                                  label: 'Porcentaje de quórum',
                                  controller: _quorumPercentageCtrl,
                                  keyboardType:
                                      const TextInputType.numberWithOptions(
                                    decimal: true,
                                  ),
                                ),
                                field(
                                  label: 'Fecha de inicio',
                                  controller: _startDateCtrl,
                                  keyboardType: TextInputType.datetime,
                                ),
                                field(
                                  label: 'Fecha de fin',
                                  controller: _endDateCtrl,
                                  keyboardType: TextInputType.datetime,
                                ),
                              ],
                            ),
                            const SizedBox(height: 8),
                            SwitchListTile.adaptive(
                              value: _requiresQuorum,
                              onChanged: _saving
                                  ? null
                                  : (value) {
                                      setState(() {
                                        _requiresQuorum = value;
                                      });
                                    },
                              title: const Text('Requiere quórum'),
                            ),
                            if (_errorMessage != null) ...[
                              const SizedBox(height: 16),
                              Container(
                                width: double.infinity,
                                padding: const EdgeInsets.all(12),
                                decoration: BoxDecoration(
                                  color: cs.errorContainer,
                                  borderRadius: BorderRadius.circular(14),
                                ),
                                child: Row(
                                  crossAxisAlignment: CrossAxisAlignment.start,
                                  children: [
                                    Icon(Icons.error_outline, color: cs.error),
                                    const SizedBox(width: 12),
                                    Expanded(
                                      child: Text(
                                        _errorMessage!,
                                        style: tt.bodyMedium?.copyWith(
                                          color: cs.onErrorContainer,
                                          fontWeight: FontWeight.w600,
                                        ),
                                      ),
                                    ),
                                  ],
                                ),
                              ),
                            ],
                            const SizedBox(height: 24),
                            Row(
                              children: [
                                Expanded(
                                  child: OutlinedButton(
                                    onPressed: _saving
                                        ? null
                                        : () => Navigator.of(context).pop(),
                                    child: const Text('Cancelar'),
                                  ),
                                ),
                                const SizedBox(width: 12),
                                Expanded(
                                  child: FilledButton.icon(
                                    onPressed: _saving ? null : _handleSubmit,
                                    icon: _saving
                                        ? SizedBox(
                                            width: 18,
                                            height: 18,
                                            child: CircularProgressIndicator(
                                              strokeWidth: 2.2,
                                              color: cs.onPrimary,
                                            ),
                                          )
                                        : const Icon(Icons.save_outlined),
                                    label: Text(
                                      _saving
                                          ? 'Guardando...'
                                          : widget.actionLabel,
                                    ),
                                  ),
                                ),
                              ],
                            ),
                          ],
                        ),
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}

class _AddVotingGroupFab
    extends GetWidget<HorizontalPropertyVotingController> {
  final String controllerTag;
  const _AddVotingGroupFab({required this.controllerTag});

  @override
  String? get tag => controllerTag;

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    final tt = Theme.of(context).textTheme;
    return Obx(() {
      final isBusy = controller.isMutationBusy.value;
      final gradientColors = isBusy
          ? [
              cs.primary.withValues(alpha: .45),
              cs.secondary.withValues(alpha: .45),
            ]
          : [cs.primary, cs.secondary];
      final shadowColor = cs.primary.withValues(alpha: isBusy ? .18 : .32);

      final labelStyle = tt.labelLarge?.copyWith(
        fontWeight: FontWeight.w800,
        color: cs.onPrimary,
      );

      final label = isBusy
          ? Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                SizedBox(
                  width: 18,
                  height: 18,
                  child: CircularProgressIndicator(
                    strokeWidth: 2.2,
                    color: cs.onPrimary,
                  ),
                ),
                const SizedBox(width: 12),
                Text('Procesando...', style: labelStyle),
              ],
            )
          : Text('Agregar grupo', style: labelStyle);

      return Container(
        decoration: BoxDecoration(
          gradient: LinearGradient(
            colors: gradientColors,
            begin: Alignment.topLeft,
            end: Alignment.bottomRight,
          ),
          borderRadius: BorderRadius.circular(28),
          boxShadow: [
            BoxShadow(
              color: shadowColor,
              blurRadius: 18,
              offset: const Offset(0, 12),
            ),
          ],
        ),
        child: FloatingActionButton.extended(
          heroTag: 'add-voting-group-$controllerTag',
          elevation: 0,
          backgroundColor: Colors.transparent,
          foregroundColor: cs.onPrimary,
          onPressed:
              isBusy ? null : () => _openCreateSheet(context, controllerTag),
          icon: isBusy ? null : const Icon(Icons.add, size: 24),
          label: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 6),
            child: label,
          ),
        ),
      );
    });
  }
}
